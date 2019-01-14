package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/rc452860/vnet/config"
	"github.com/rc452860/vnet/log"
	"github.com/rc452860/vnet/proxy"
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

	if conf.DbConfig.Rate == 0 {
		survey.AskOne(&survey.Input{
			Message: "what is your want rate?",
			Default: "1",
		}, &conf.DbConfig.Rate, nil)
	}

	if conf.DbConfig.SyncTime == 0 {
		survey.AskOne(&survey.Input{
			Message: "what is your want database sync time?",
			Default: "3",
		}, &conf.DbConfig.Rate, nil)
	}
	// start database traffic monitor
	go DBTrafficMonitor(ctx)

	go func(ctx context.Context) {
		tick := time.Tick(3 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			config.SaveConfig()
			users, err := GetEnableUser()
			if err != nil {
				log.Err(err)
				return
			}

			sslist := service.CurrentShadowsocksService().List()
			userMap := make(map[int]User)
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

			for _, ss := range sslist {
				flag := false
				for _, user := range users {
					if ss.Port == user.Port {
						flag = true
						break
					}
				}
				if !flag {
					log.Info("stop port [%d]", ss.Port)
					service.CurrentShadowsocksService().Del(ss.Port)
				}
			}
			UpdateTrafficByUser(userMap)
			//TODO update user upload and download
			<-tick
		}
	}(ctx)
}

type DBTraffic struct {
	Port int
	Up   uint64
	Down uint64
}

var trafficTable = make(map[int]*DBTraffic)

func UpdateTrafficByUser(users map[int]User) {
	// TODO move to log
	defer func() {
		var err error
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
			// invalidate rep
			log.Err(err)
		}
	}()

	whenUp := new(bytes.Buffer)
	whenDown := new(bytes.Buffer)
	whenPort := new(bytes.Buffer)
	insertBuf := new(bytes.Buffer)
	conf := config.CurrentConfig().DbConfig
	for _, v := range trafficTable {
		// if the traffic less then 64 kb it will no necessary to update
		if v.Up+v.Down < 512*1024 {
			continue
		}
		whenUp.WriteString(fmt.Sprintf(" WHEN %v THEN u+%v", v.Port, v.Up))
		whenDown.WriteString(fmt.Sprintf(" WHEN %v THEN d+%v", v.Port, v.Down))
		whenPort.WriteString(fmt.Sprintf("%v,", v.Port))
		traffic, err := datasize.HumanSize(v.Up + v.Down)
		if err != nil {
			log.Error("traffic parse error")
			traffic = "0"
		}
		insertBuf.WriteString(fmt.Sprintf("(NULL,%v,%v,%v,%v,%f,'%s',unix_timestamp()),",
			users[v.Port].Id,
			v.Up,
			v.Down,
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
	query := fmt.Sprintf("UPDATE user SET u = CASE port%s END,"+
		"d = CASE port%s END,t = unix_timestamp() WHERE port IN (%s)",
		whenUp.String(),
		whenDown.String(),

		portStr[:len(portStr)-1])
	inserStr := insertBuf.String()
	insert := fmt.Sprintf("INSERT INTO `user_traffic_log` (`id`, `user_id`, `u`, `d`, "+
		"`node_id`, `rate`, `traffic`, `log_time`) VALUES %s", inserStr[:len(inserStr)-1])
	connect, err := Connect()
	if err != nil {
		log.Err(err)
		return
	}
	connect.Exec(query)
	connect.Exec(insert)
}

func DBTrafficMonitor(ctx context.Context) {
	traffic := make(chan proxy.TrafficMessage, 128)
	proxy.RegisterTrafficHandle(traffic)
	var data proxy.TrafficMessage
	// count traffic
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data = <-traffic:
			}
			_, port, err := net.SplitHostPort(data.LAddr)
			if err != nil {
				log.Err(err)
				continue
			}
			porti, err := strconv.Atoi(port)
			if err != nil {
				log.Err(err)
				continue
			}
			if trafficTable[porti] == nil {
				trafficTable[porti] = &DBTraffic{
					Port: porti,
					Up:   data.UpBytes,
					Down: data.DownBytes,
				}
			} else {
				trafficTable[porti].Up += data.UpBytes
				trafficTable[porti].Down += data.DownBytes
			}
		}
	}(ctx)
}

func StartShadowsocks(user User) {
	limit, err := datasize.HumanSize(user.Limit)
	if err != nil {
		log.Error("limit: %d, error info:%s", user.Limit, err.Error())
		return
	}
	err = service.CurrentShadowsocksService().Add("0.0.0.0",
		user.Method,
		user.Password,
		user.Port,
		limit,
		3*time.Second)
	if err != nil {
		log.Info("[%d] add failure, case %s", user.Port, err.Error())
	}
	err = service.CurrentShadowsocksService().Start(user.Port)
	if err != nil {
		log.Info("[%d] started failure, case: %s", user.Port, err.Error())
	}
}
