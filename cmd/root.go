/*
 * Copyright (C) 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bzp2010/webvpn/internal/handler"
	"github.com/bzp2010/webvpn/internal/model"
	"github.com/bzp2010/webvpn/internal/utils"
)

var configFile string
var rootCmd = &cobra.Command{
	Use:   "webvpn",
	Short: "WebVPN is a zero-trust gateway for proxy application",
	Long:  `A flexible and configurable zero-trust gateway, it provides pluggable authentication and authorization for applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer(true, false)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Log().Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onInitialize)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "config file")
}

func onInitialize() {
	// load config
	initConfig()

	// read proxy policy
	initPolicy()
}

func initConfig() {
	// set config file path
	viper.SetConfigFile(configFile)

	// reading environments
	viper.AutomaticEnv()

	// reading config file
	if err := viper.ReadInConfig(); err != nil {
		utils.Log().Error(err)
	}

	// watch and auto reload config file
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.ReadInConfig(); err != nil {
			utils.Log().Error(err)
		}
		utils.Log().Debugf("config auto reload: %s", in.Name)
	})

	// setting debug
	utils.Log().SetDebug(viper.GetBool("debug"))
}

func initPolicy() {
	var policySlice []*model.Policy
	err := viper.UnmarshalKey("policy", &policySlice)
	if err != nil {
		utils.Log().Error(err)
		return
	}

	for _, v := range policySlice {
		handler.PolicyMap.Set(v.From, v)
	}

	utils.Log().Debug("policyMap:", handler.PolicyMap)
}
