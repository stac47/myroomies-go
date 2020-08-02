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
		"MongoDB URL (e.g. mongodb://localhost:27017). If no value "+
			"or an invalid MongoDB URL is provided, MyRoomies will store the data "+
			"in memory and will be lost at server shutdown.")
	rootCmd.Flags().StringVarP(&serverConfig.BindTo, "bind-to", "", ":8080",
		"Bind MyRoomies server to a given address")
	rootCmd.Flags().StringVarP(&serverConfig.CertificatePath, "cert-file", "", "",
		"Path to the certificate to enable TLS connections. To be used with --key-file.")
	rootCmd.Flags().StringVarP(&serverConfig.KeyPath, "key-file", "", "",
		"Path to the private key to enable TLS connections. To be used with --cert-file.")
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
	return rest.Start(serverConfig)
}

func main() {
	log.Debug(strings.Join(os.Args, " "))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
