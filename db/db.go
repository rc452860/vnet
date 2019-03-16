package db

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rc452860/vnet/proxy/server"

	"github.com/rc452860/vnet/utils/addr"

	"github.com/rc452860/vnet/record"

	"github.com/AlecAivazis/survey"
	"github.com/rc452860/vnet/common/config"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/service"
	"github.com/rc452860/vnet/utils/datasize"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id             int64  `gorm:"column:id;primary_key"`
	Port           int    `gorm:"column:port"`
	Enable         bool   `gorm:"column:enable"`
	Update         int64  `gorm:"column:u"`
	Download       int64  `gorm:"column:d"`
	TransferEnable int64  `gorm:"column:transfer_enable"`
	Method         string `gorm:"column:method"`
	Password       string `gorm:"column:passwd"`
	Limit          uint64 `gorm:"column:speed_limit_per_user"`
}

type DBTraffic struct {
	Port         int
	Up           uint64
	Down         uint64
	ConnectCount int
	Connects     map[string]bool
}

var (
	trafficTable = make(map[int]*DBTraffic)
	startTime    = time.Now()
)

func (this User) TableName() string {
	return "user"
}

func Connect() (*gorm.DB, error) {

	var (
		username string
		password string
		database string
		host     string
		port     string
	)
	pattern := "username:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local"

	c := config.CurrentConfig()
	if c == nil {
		return nil, fmt.Errorf("config dbconfig is nil")
	}
	// log.Info("dbConfig: %s", c.DbConfig)
	username = c.DbConfig.User
	password = c.DbConfig.Passwd
	database = c.DbConfig.Database
	host = c.DbConfig.Host
	port = c.DbConfig.Port
	replacer := strings.NewReplacer(
		"username",
		username,
		"password",
		password,
		"host",
		host,
		"port",
		port,
		"database",
		database,
	)
	conntionStr := replacer.Replace(pattern)
	// log.Info("connection string: %s", conntionStr)
	db, err := gorm.Open("mysql", conntionStr)
	return db, err
}

func GetEnableUser() ([]User, error) {
	connect, err := Connect()
	if err != nil {
		return nil, err
	}
	defer connect.Close()

	var userList []User
	connect.Where("port != 0 AND enable = 1").Find(&userList)
	return userList, nil
}

func GetUserWithLevel(level int) ([]User, error) {
	connect, err := Connect()
	if err != nil {
		return nil, err
	}
	defer connect.Close()

	var userList []User
	connect.Where("port != 0 AND enable = 1 AND level >= level").Find(&userList)
	return userList, nil
}

func DbStarted(ctx context.Context) {
	conf := config.CurrentConfig()
	if conf.DbConfig.Host == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your host address?",
		}, &conf.DbConfig.Host, nil)
	}

	if conf.DbConfig.Port == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your database port?",
			Default: "3306",
		}, &conf.DbConfig.Port, nil)
	}

	if conf.DbConfig.User == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your username?",
		}, &conf.DbConfig.User, nil)
	}

	if conf.DbConfig.Passwd == "" {
		survey.AskOne(&survey.Password{
			Message: "what is your password?",
		}, &conf.DbConfig.Passwd, nil)
	}

	if conf.DbConfig.Database == "" {
		survey.AskOne(&survey.Input{
			Message: "what is your database name?",
		}, &conf.DbConfig.Database, nil)
	}

	if conf.DbConfig.NodeId == 0 {
		survey.AskOne(&survey.Input{
			Message: "what is your node id?",
		}, &conf.DbConfig.NodeId, nil)
	}

	if conf.DbConfig.Rate == -1 {
		survey.AskOne(&survey.Input{
			Message: "what is your want rate?",
			Default: "1",
		}, &conf.DbConfig.Rate, nil)
	}

	if conf.DbConfig.SyncTime == 0 {
		survey.AskOne(&survey.Input{
			Message: "what is your want database sync time (unit: Millisencod)?",
			Default: "60000",
		}, &conf.DbConfig.SyncTime, nil)
	}

	if conf.DbConfig.OnlineSyncTime == 0 {
		survey.AskOne(&survey.Input{
			Message: "what is the value your want set online user sync time (unit: Millisencod)?",
			Default: "60000",
		}, &conf.DbConfig.OnlineSyncTime, nil)
	}
	// save config
	config.SaveConfig()

	// start database traffic monitor
	go DBTrafficMonitor(ctx)
	// start database service monitor
	go DBServiceMonitor(ctx)
	// online monitor
	go DBOnlineMonitor(ctx)
}

// DBOnlineMonitor auto upload last one minute ips
func DBOnlineMonitor(ctx context.Context) {
	conf := config.CurrentConfig().DbConfig
	tick := time.Tick(time.Duration(conf.OnlineSyncTime) * time.Millisecond)
	grmInstance := record.GetGRMInstanceWithTick(time.Duration(conf.OnlineSyncTime) * time.Millisecond)
	connect, err := Connect()
	if err != nil {
		log.Err(err)
		return
	}
	defer connect.Close()

	for {
		ss_node_info := fmt.Sprintf("INSERT INTO `ss_node_info` (`id`,`node_id`,`uptime`,`load`,`log_time`)"+
			"VALUES (NULL,%v,%v,'%s',unix_timestamp())", conf.NodeId, time.Since(startTime).Seconds(), formatLoad())
		ss_node_online_log := fmt.Sprintf("INSERT INTO `ss_node_online_log` (id,node_id,online_user,log_time) "+
			" VALUES (NULL,%v,%v,unix_timestamp())", conf.NodeId, grmInstance.GetLastOneMinuteOnlineCount())
		ss_node_ip := bytes.NewBufferString("")
		ss_node_ip.WriteString("INSERT INTO `ss_node_ip`(id,node_id,port,type,ip,created_at) VALUES")
		headLen := ss_node_ip.Len()
		userOnline := grmInstance.GetLastOneMinuteOnlineByPort()
		for k, v := range userOnline {
			tcpIps := bytes.NewBufferString("")
			udpIps := bytes.NewBufferString("")

			for _, ip := range v {
				if strings.Contains(ip.Network(), "tcp") {
					tcpIps.WriteString(addr.GetIPFromAddr(ip))
					tcpIps.WriteString(",")
				}
				if strings.Contains(ip.Network(), "tcp") {
					udpIps.WriteString(addr.GetIPFromAddr(ip))
					udpIps.WriteString(",")
				}
			}
			if tcpIps.Len() > 1 {
				ss_node_ip.WriteString(fmt.Sprintf("(NULL,%v,%v,'%s','%s',unix_timestamp()),", conf.NodeId, k, "tcp", tcpIps.String()[:tcpIps.Len()-1]))

			}
			if udpIps.Len() > 1 {
				ss_node_ip.WriteString(fmt.Sprintf("(NULL,%v,%v,'%s','%s',unix_timestamp()),", conf.NodeId, k, "udp", udpIps.String()[:udpIps.Len()-1]))
			}
		}
		connect.Exec(ss_node_online_log)
		connect.Exec(ss_node_info)
		if ss_node_ip.Len() > headLen {
			connect.Exec(ss_node_ip.String()[:ss_node_ip.Len()-1])
		}

		<-tick
	}
}

// DBServiceMonitor is shadowsocks service monitor
// it will keep watch the database config change,and period
// apply the config to our program
func DBServiceMonitor(ctx context.Context) {
	conf := config.CurrentConfig().DbConfig
	tick := time.Tick(time.Duration(conf.SyncTime) * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var (
			users []User
			err   error
		)
		if conf.Level != -1 {
			users, err = GetEnableUser()
		} else {
			users, err = GetUserWithLevel(conf.Level)
		}
		if err != nil {
			log.Err(err)
			continue
		}

		sslist := service.CurrentShadowsocksService().List()
		userMap := make(map[int]User)

		for _, ss := range sslist {
			flag := false
			for _, user := range users {
				if ss.Port == user.Port && ss.Method == user.Method && ss.Password == user.Password {
					flag = true
					break
				}
			}
			if !flag {
				log.Info("stop port [%d]", ss.Port)
				service.CurrentShadowsocksService().Del(ss.Port)
			}
		}

		for _, user := range users {
			userMap[user.Port] = user
			flag := false
			for _, ss := range sslist {
				if ss.Port == user.Port {
					flag = true
					break
				}
			}
			if !flag {
				log.Info("start port [%d]", user.Port)
				StartShadowsocks(user)
			}
		}

		UpdateTrafficByUser(userMap)
		<-tick
	}
}

func UpdateTrafficByUser(users map[int]User) {
	// TODO move to log
	// defer func() {
	// 	var err error
	// 	if r := recover(); r != nil {
	// 		// find out exactly what the error was and set err
	// 		switch x := r.(type) {
	// 		case string:
	// 			err = errors.New(x)
	// 		case error:
	// 			err = x
	// 		default:
	// 			// Fallback err (per specs, error strings should be lowercase w/o punctuation
	// 			err = errors.New("unknown panic")
	// 		}
	// 		// invalidate rep
	// 		log.Err(err)
	// 	}
	// }()

	whenUp := new(bytes.Buffer)
	whenDown := new(bytes.Buffer)
	whenPort := new(bytes.Buffer)
	insertBuf := new(bytes.Buffer)
	conf := config.CurrentConfig().DbConfig
	connect, err := Connect()
	if err != nil {
		log.Err(err)
		return
	}
	defer connect.Close()
	for _, v := range trafficTable {
		// if the traffic less then 64 kb it will no necessary to update
		if v.Up+v.Down < 64*1024 {
			continue
		}
		whenUp.WriteString(fmt.Sprintf(" WHEN %v THEN u+%v", v.Port, float32(v.Down)*conf.Rate))
		whenDown.WriteString(fmt.Sprintf(" WHEN %v THEN d+%v", v.Port, float32(v.Up)*conf.Rate))
		whenPort.WriteString(fmt.Sprintf("%v,", v.Port))
		traffic, err := datasize.HumanSize((v.Up + v.Down) * uint64(conf.Rate))
		if err != nil {
			log.Error("traffic parse error")
			traffic = "0"
		}
		insertBuf.WriteString(fmt.Sprintf("(NULL,%v,%v,%v,%v,%f,'%s',unix_timestamp()),",
			users[v.Port].Id,
			v.Down,
			v.Up,
			conf.NodeId,
			conf.Rate,
			traffic))

		// assumble query done then clear traffic record and commit to database
		v.Down = 0
		v.Up = 0

	}
	if whenPort.Len() < 1 {
		return
	}
	// assembly query
	portStr := whenPort.String()
	user := fmt.Sprintf("UPDATE user SET u = CASE port%s END,"+
		"d = CASE port%s END,t = unix_timestamp() WHERE port IN (%s)",
		whenUp.String(),
		whenDown.String(),
		portStr[:len(portStr)-1])

	inserStr := insertBuf.String()
	userTrafficLog := fmt.Sprintf("INSERT INTO `user_traffic_log` (`id`, `user_id`, `u`, `d`, "+
		"`node_id`, `rate`, `traffic`, `log_time`) VALUES %s", inserStr[:len(inserStr)-1])

	connect.Exec(user)
	connect.Exec(userTrafficLog)
}

func formatLoad() string {
	return fmt.Sprintf("cpu:%v%% mem:%v%% disk:%v%%", service.GetCPUUsage(), service.GetMemUsage(), service.GetDiskUsage())
}

// DBTrafficMonitor is should using goroutine start
func DBTrafficMonitor(ctx context.Context) {
	traffic := make(chan record.Traffic, 32)
	server.RegisterTrafficHandle(traffic)
	var data record.Traffic
	// count traffic
	for {
		select {
		case <-ctx.Done():
			return
		case data = <-traffic:
		}
		port := addr.GetPortFromAddr(data.ProxyAddr)
		if trafficTable[port] == nil {
			trafficTable[port] = &DBTraffic{
				Port: port,
				Up:   data.Up,
				Down: data.Down,
			}
		} else {
			trafficTable[port].Up += data.Up
			trafficTable[port].Down += data.Down
		}
	}
}

func StartShadowsocks(user User) {
	sscon := config.CurrentConfig().ShadowsocksOptions

	err := service.CurrentShadowsocksService().Add("0.0.0.0",
		user.Method,
		user.Password,
		user.Port,
		server.ShadowsocksArgs{
			Limit:          user.Limit,
			ConnectTimeout: time.Duration(sscon.ConnectTimeout) * time.Millisecond,
			TCPSwitch:      sscon.TCPSwitch,
			UDPSwitch:      sscon.UDPSwitch,
		})
	if err != nil {
		log.Info("[%d] add failure, case %s", user.Port, err.Error())
	}
	err = service.CurrentShadowsocksService().Start(user.Port)
	if err != nil {
		log.Info("[%d] started failure, case: %s", user.Port, err.Error())
	}
}
