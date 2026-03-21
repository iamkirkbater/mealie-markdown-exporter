package cmd

import (
	"fmt"

	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/apitoken"
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/outputdirectory"
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/markdown"
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/mealie"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export recipes from Mealie as markdown files",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		if viper.GetString("base-url") == "" {
			return fmt.Errorf("base-url is required (set via --base-url flag or MME_BASE_URL env var)")
		}
		if err := resolveOutputDir(); err != nil {
			return err
		}
		err := resolveAPIToken(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := viper.GetString("base-url")
		apiToken := viper.GetString("api-token")

		log.Info("Exporting from: ", baseURL)

		client := mealie.NewClient(baseURL, apiToken)
		recipes, err := client.GetAllRecipes()
		if err != nil {
			return err
		}

		log.Infof("Retrieved %d recipes", len(recipes))

		outputDir := viper.GetString("output-dir")
		if err := markdown.WriteRecipes(afero.NewOsFs(), outputDir, recipes); err != nil {
			return err
		}

		log.Infof("Exported %d recipes to %s", len(recipes), outputDir)
		return nil
	},
}

func resolveOutputDir() error {
	return outputdirectory.Resolve(afero.NewOsFs(), viper.GetString("output-dir"))
}

func resolveAPIToken(cmd *cobra.Command) error {
	rawToken := viper.GetString("api-token")
	if rawToken == "" {
		return fmt.Errorf("api-token is required (set via --api-token flag or MME_API_TOKEN env var)")
	}
	token, err := apitoken.Resolve(afero.NewOsFs(), rawToken)
	if err != nil {
		return err
	}
	viper.Set("api-token", token)
	if cmd.Flags().Changed("api-token") && token == rawToken {
		log.Warn("API token provided directly via flag; this is sensitive and will be stored in your shell history, and may be exposed to other users.")
		log.Warn("Consider using the MME_API_TOKEN env var instead, or storing the token in a file and using --api-token file:///path/to/token.")
	}
	return nil
}

func init() {
	exportCmd.Flags().String("base-url", "", "Base URL of the Mealie instance")
	exportCmd.Flags().String("api-token", "", "API token for the Mealie instance")
	exportCmd.Flags().String("output-dir", outputdirectory.DefaultOutputDir, "Directory to export markdown files to")
	rootCmd.AddCommand(exportCmd)
}
