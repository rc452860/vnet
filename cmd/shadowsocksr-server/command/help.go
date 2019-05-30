package command

// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rc452860/vnet/common/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "shadowoscksr webapi version",
}


func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().String("config", "config.json", "config file default: config.json")
	_ = viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	for _, item := range flagConfigs {
		rootCmd.Flags().String(item.Name, item.Default, item.Usage)
		_ = viper.BindPFlag(item.Name, rootCmd.Flags().Lookup(item.Name))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(viper.GetString("config"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute(fn func()) {
	rootCmd.Run = runWrap(fn)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkRequired() bool {
	for _, item := range flagConfigs {
		if item.Required && viper.GetString(item.Name) == "" {
			log.Warn("miss param:" + item.Name)
			return false
		}
	}
	return true
}

func runWrap(fn func()) func(cmd *cobra.Command, args []string){
	return func(cmd *cobra.Command, args []string){
		if !checkRequired() {
			_ = cmd.Help()
			return
		}
		fn()
	}
}

