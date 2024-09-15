// Package cmd /*
package cmd

import (
	"os"

	"log"
	"yuka/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/go-playground/validator/v10"
)

/*
Example commands to be built out
yukactl user create
yukactl user root
yukactl user delete

yukactl organization create
yukactl organization root
yukactl organization delete

yukactl application create
yukactl application root
yukactl application delete
*/

type rootOptions struct {
	ApiserverAddress string `flag:"apiserver-address" validate:"required"`
}

var _rootOptions rootOptions

var validationFns = map[string]func(validator.FieldLevel) bool{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yukactl",
	Short: "A useful client for dealing with a useful network!",
	Long: `TBA.
	Run "yukactl help" for more information.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateAndUnmarshal(cmd, &_rootOptions, validationFns); err != nil {
			log.Fatalln(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := utils.GetLogger()
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Sugar().Debug("Application root called")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringP("apiserver-address", "a", "localhost:8080", "Address of the yuka api server.")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
