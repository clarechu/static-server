package router

import (
	"github.com/spf13/cobra"
	"net/http"
)

type Root struct {
	Port    int32  `json:"port"`
	FileDir string `json:"file_dir"`
	Path    string `json:"path"`
	Index   string `json:"index"`
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
	cmd.PersistentFlags().Int32VarP(&args.Port, "port", "p", 8080, "static file server ports")
	cmd.PersistentFlags().StringVarP(&args.Path, "path", "P", "/console", "static file server path")
	cmd.PersistentFlags().StringVarP(&args.FileDir, "file", "f", "./dist", "static file path")
	cmd.PersistentFlags().StringVarP(&args.Index, "index", "i", "./dist/index.html", "static file path index.html")
}
