package config

import (
	"fmt"
	"testing"
)

func Test_config(t *testing.T) {
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Print(err)
	}
	// config.DbConfig = &DbConfig{
	// 	Host: "mysql.vnet.club",
	// }
	// config.ShadowsocksOptions = &ShadowsocksOptions{
	// 	TcpTimeout: 12,
	// }
	config.ShadowsocksOptions.ConnectTimeout = 3
	// fmt.Print(config)
	err = SaveConfig()
	if err != nil {
		fmt.Print(err)
	}
}

func Benchmark_config(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		CurrentConfig()
	}
}
