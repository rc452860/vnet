package utils

import (
	"fmt"
	"testing"
)

func Test_Config(t *testing.T) {
	config, err := ConfigFactory("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	config.Map.PutString("sakura", "killer")
	config.Map.PutInt("abc", 1)
	subMap := ConfigMap{}
	subMap.PutString("funck", "avc")
	config.Map.PutConfigMap("subMap", subMap)
	config.WriteConfig()
}

func Test_ConfigRead(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	config, err := ConfigFactory("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.Map.GetInt("abc"))
	fmt.Println(config.Map.GetInt("bcd"))
	fmt.Println(config.Map.GetString("bcd"))
	fmt.Println(config.Map.GetConfigMap("bcd"))
	fmt.Println(config.Map.GetString("sakura"))
	fmt.Println(config.Map.GetConfigMap("subMap").GetString("funck"))
}
