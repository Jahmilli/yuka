package client

import (
	"log"
	"yuka/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

type clientOptions struct {
	ApiserverAddress string `flag:"apiserver-address" validate:"required"`
}

var _clientOptions clientOptions

var validationFns = map[string]func(validator.FieldLevel) bool{}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Starts up a client",
	Long:  `Starts up a client. Run "yuka client --help" for more information.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateAndUnmarshal(cmd, &_clientOptions, validationFns); err != nil {
			log.Fatalln(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func SubCommand() *cobra.Command {
	clientCmd.PersistentFlags().StringP("apiserver-address", "a", "localhost:8080", "Address of the yuka api server.")
	return clientCmd
}
