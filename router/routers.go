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
		Long:  `http-server 加载静态资源.`,
		Run: func(cmd *cobra.Command, args []string) {
			server := NewServer(ag)
			server.Run()
		},
	}
	addFlag(rootCmd, ag)
	return rootCmd
}

func addFlag(cmd *cobra.Command, args *Root) {
	cmd.PersistentFlags().Int32Var(&args.Port, "port", 8080, "static file server ports")
	cmd.PersistentFlags().StringVarP(&args.Path, "basicPath", "p", "/", "url root path")
	cmd.PersistentFlags().StringVarP(&args.FileDir, "file", "f", "./dist", "static file path")
	cmd.PersistentFlags().StringVar(&args.PublicPath, "publicPath", "/", "The base URL your application bundle will be deployed")
}
