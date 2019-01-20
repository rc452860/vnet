package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rc452860/vnet/comm/log"
	"github.com/rc452860/vnet/utils"
)

var (
	config     *Config
	configFile string
)

// Config is global config
type Config struct {
	Mode               string             `json:"mode"`
	DbConfig           DbConfig           `json:"dbconfig"`
	ShadowsocksOptions ShadowsocksOptions `json:"shadowsocks_options"`
}

// DbConfig is global database config
type DbConfig struct {
	Host           string  `json:"host"`
	User           string  `json:"user"`
	Passwd         string  `json:"passwd"`
	Port           string  `json:"port"`
	Database       string  `json:"database"`
	Rate           float32 `json:"rate"`
	NodeId         int     `json:"node_id`
	SyncTime       int     `json:'sync_time'`
	OnlineSyncTime int     `json:'online_sync_time'`
}

type ShadowsocksOptions struct {
	TCPTimeout int `json:"tcp_timeout"`
	UDPTimeout int `json:"udp_timeout"`
}

func CurrentConfig() *Config {
	if config == nil {
		conf, err := LoadDefault()
		if err != nil {
			panic(err)
		}
		config = conf
	}
	return config
}

func LoadDefault() (*Config, error) {
	return LoadConfig("config.json")
}

func LoadConfig(file string) (*Config, error) {
	utils.RLock(file)
	defer utils.RUnLock(file)
	if !utils.IsFileExist(file) {
		absFile, err := filepath.Abs(file)
		if err != nil {
			log.Err(err)
		} else {
			log.Warn("%s is not exist", absFile)
		}
		configFile = file
		config = &Config{
			Mode:     "db",
			DbConfig: DbConfig{},
		}
		data, _ := json.MarshalIndent(config, "", "    ")
		ioutil.WriteFile(configFile, data, 0644)
		return config, nil
	}
	config = &Config{
		Mode: "bare",
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("resolve config file failed: %v", err)
	}
	configFile = file
	return config, nil
}

func SaveConfig() error {
	if config == nil {
		return fmt.Errorf("not config loaded!")
	}

	data, err := json.MarshalIndent(config, "", "    ")

	if err != nil {
		return fmt.Errorf("config marshal failed!")
	}

	return ioutil.WriteFile(configFile, data, 0644)
}

func (self Config) String() string {
	data, err := json.MarshalIndent(self, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(data)
}
