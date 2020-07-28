package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/stac47/myroomies/internal/server/rest"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringVarP(&serverConfig.Storage, "storage", "", "",
		"MongoDB URI (e.g. mongodb://localhost:27017). By default, the "+
			"data will be stored in memory and will be lost when the server "+
			"shutdown.")
	rootCmd.Flags().StringVarP(&serverConfig.BindTo, "bind-to", "", ":8080",
		"Bind MyRoomies server to a given address")
}

var (
	rootCmd = &cobra.Command{
		Use:   "myroomies-server",
		Short: "The server part of MyRoomies",
		Long: `MyRoomies is a tool that helps people leaving in houseshares to
organize their tasks and expenses`,
		RunE: startServer,
	}

	serverConfig rest.ServerConfig
)

func startServer(cmd *cobra.Command, args []string) error {
	rest.Start(serverConfig)
	return nil
}

func main() {
	log.Println(strings.Join(os.Args, " "))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
