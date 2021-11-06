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
	"sync"

	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/core"
	"github.com/bzp2010/webvpn/internal/utils"
)

// root command of run server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run WebVPN server",
}

func init() {
	// run only public server
	serveCmd.AddCommand(&cobra.Command{
		Use: "public",
		Short: "Run public server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(true, false)
		},
	})

	// run only admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "admin",
		Short: "Run admin server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(false, true)
		},
	})

	// run public and admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "all",
		Short: "Run all server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(true, true)
		},
	})

	rootCmd.AddCommand(serveCmd)
}

func startServer(isStartPublic, isStartAdmin bool)  {
	// initialize webvpn server
	s, err := core.Server(&core.Options{
		Public: isStartPublic,
		Admin:  isStartAdmin,
	})
	if err != nil {
		utils.Log().Errorf("webvpn server initialize failed: %s", err.Error())
		return
	}

	// start webvpn server
	s.Start()

	// create waitgroup
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}