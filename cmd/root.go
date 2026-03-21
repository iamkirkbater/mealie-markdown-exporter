package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mealie-markdown-exporter",
	Short: "Export Mealie recipes to markdown files",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set any log level if it's been configured.
		level := viper.GetString("log-level")
		logLevel, err := log.ParseLevel(level)
		if err != nil {
			return fmt.Errorf("invalid log level %q: must be one of debug, info, warn, error", level)
		}
		log.SetLevel(logLevel)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("MME")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.BindPFlags(rootCmd.PersistentFlags())
}
