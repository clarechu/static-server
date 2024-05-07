// Copyright (c) 2021 The static-server Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"github.com/spf13/cobra"
	"net/http"
)

type Root struct {
	Port       int32
	FileDir    string
	Path       string
	PublicPath string
	IsGzip     bool
	Config     string
}

type Server struct {
	sv *http.Server
}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd(args []string) *cobra.Command {
	ag := &Root{}
	rootCmd := &cobra.Command{
		Use:   "http-server",
		Short: "http-server ...",
		Long: `
Tips  Find more information at: https://github.com/clarechu/static-server

Example:

# http-server dist .
http-server dist

# http-server set publicPath .
http-server ./dist --publicPath /a

# set port .
http-server dist --port 8080

# 开启页面gzip
http-server dist --port 8080 --gzip true
`,
		Run: func(cmd *cobra.Command, args []string) {
			server := NewServer(ag)
			server.Run()
		},
	}
	addFlag(rootCmd, ag)
	rootCmd.AddCommand(VersionCommand())
	return rootCmd
}

func addFlag(cmd *cobra.Command, args *Root) {
	cmd.PersistentFlags().Int32Var(&args.Port, "port", 8080, "static file server ports")
	cmd.PersistentFlags().BoolVar(&args.IsGzip, "gzip", false, "gzip  (default false)")
	cmd.PersistentFlags().StringVarP(&args.Path, "basicPath", "p", "/", "url root path")
	cmd.PersistentFlags().StringVarP(&args.FileDir, "file", "f", "./dist", "static file path")
	cmd.PersistentFlags().StringVar(&args.PublicPath, "publicPath", "/", "The base URL your application bundle will be deployed")
	cmd.PersistentFlags().StringVar(&args.Config, "config", "", "The static configuration file path")

}
