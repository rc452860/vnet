package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	File string
	Map  ConfigMap
}

type ConfigMap map[string]interface{}

var ConfigInstance map[string]*Config

func init() {
	ConfigInstance = make(map[string]*Config)
}

func ConfigFactory(file string) (config *Config) {
	if ConfigInstance[file] != nil {
		return ConfigInstance[file]
	}
	config = &Config{
		File: file,
		Map:  ConfigMap{},
	}
	if err := config.ReadConfig(); err != nil {
		panic("can not read file")
	}
	return config
}

func (this *Config) ReadConfig() error {
	if IsFileExist(this.File) {
		result, err := ioutil.ReadFile(this.File)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(result, &this.Map); err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}

func (this *Config) WriteConfig() error {
	file, err := os.OpenFile(this.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	result, err := json.MarshalIndent(this.Map, "", "    ")
	if err != nil {
		return err
	}
	w, err := file.Write(result)
	if err != nil && w < 0 {
		return err
	}
	return nil
}

func (this ConfigMap) GetInt(name string) int {
	value := this[name]
	if value == nil {
		return 0
	}
	switch value.(type) {
	case float64:
		return int(value.(float64))
	case string:
		result, err := strconv.Atoi(value.(string))
		if err != nil {
			panic(err)
		}
		return result
	default:
		return 0
	}
}

func (this ConfigMap) GetString(name string) string {
	value := this[name]
	if value == nil {
		return ""
	}
	switch value.(type) {
	case float64:
		return strconv.Itoa(int(value.(float64)))
	case string:
		return value.(string)
	default:
		return "0"
	}
}

func (this ConfigMap) GetConfigMap(name string) ConfigMap {
	value := this[name]
	if value == nil {
		return nil
	}
	return ConfigMap(this[name].(map[string]interface{}))
}

func (this ConfigMap) PutInt(key string, value int) {
	this[key] = value
}

func (this ConfigMap) PutString(key string, value string) {
	this[key] = value
}

func (this ConfigMap) PutConfigMap(key string, value ConfigMap) {
	this[key] = value
}

func (this ConfigMap) String() string {
	result, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(result)
}
