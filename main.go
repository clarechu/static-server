package main

import (
	"github.com/ClareChu/static-server/router"
	"os"
)

func main() {
	rootCmd := router.GetRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
