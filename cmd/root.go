/*
Copyright © 2024 Pete Wall <pete@petewall.net>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/petewall/event-logger/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "event-logger",
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug mode enabled")
		}

		path := viper.GetString("file")
		if path == "" {
			return fmt.Errorf("missing file path")
		}

		events := &internal.FilesystemDatastore{
			Path: path,
		}
		err := events.Initialize()
		if err != nil {
			return fmt.Errorf("unable to initialize the event list: %w", err)
		}

		server := &internal.Server{
			Events: events,
			Port:   viper.GetInt("port"),
		}
		return server.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	rootCmd.Flags().IntP("port", "p", internal.DefaultPort, "Port number to listen on (env: PORT)")
	_ = viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().StringP("file", "f", "events.log", "File to write events to (env: FILE)")
	_ = viper.BindPFlag("file", rootCmd.Flags().Lookup("file"))

	rootCmd.Flags().Bool("debug", false, "Enable debug logging (env: DEBUG)")
	_ = viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
	viper.AutomaticEnv()
}
