package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kaato137/quickrest/internal"
	"github.com/kaato137/quickrest/internal/conf"
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "quickrest",
	Short: "QuickREST is the quick way to mock API",
	Long:  `QuickREST is a convenient tool for quickly mocking API endpoints when you don't have them readily available.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := conf.LoadConfigFromFile(configPath)
		if err != nil {
			if errors.Is(err, conf.ErrDefaultPathNotFound) {
				cmd.PrintErr(ErrTxtDefaultPathNotFound)
			}
			return err
		}

		server, err := internal.NewServerFromConfig(cfg)
		cobra.CheckErr(err)

		defer server.Close()

		http.ListenAndServe(cfg.Address, server)
		cobra.CheckErr(err)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to a configuration file")

	rootCmd.AddCommand(generateDefaultConfigCmd)
}

func Execute(version, build string) error {
	rootCmd.Version = fmt.Sprintf("%s.%s", version, build)
	return rootCmd.Execute()
}

const ErrTxtDefaultPathNotFound = `The expected default configuration file is named 'quickrest.yml' or 'quickrest.yaml'.
Alternatively, you can specify a custom configuration file using:
	quickrest -c your_custom_config.yml

You can also generate the default configuration file by running:
	quickrest generate-default-config
or
	quickrest gdc
`
