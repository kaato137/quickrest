package cmd

import (
	"github.com/kaato137/quickrest/internal/conf"
	"github.com/spf13/cobra"
)

var generateDefaultConfigCmd = &cobra.Command{
	Use:     "generate-default-config",
	Aliases: []string{"gdc"},
	Short:   "Generate default configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return conf.GenerateDefault()
	},
}
