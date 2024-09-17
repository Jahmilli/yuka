package client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"yuka/internal/client"
	"yuka/pkg/utils"

	"github.com/spf13/cobra"
)

type startOptions struct {
	Hostname string `flag:"hostname"`
}

var _startOptions startOptions

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts up yukactl",
	Long: `Starts a connection from yukactl to the server using default configuration. 
Run "yukactl client start --help" for more information.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateAndUnmarshal(cmd, &_startOptions, validationFns); err != nil {
			log.Fatalln(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		logger, err := utils.GetLogger()
		if err != nil {
			logger.Fatal(err.Error())
		}

		apiserverAddress, _ := cmd.Flags().GetString("apiserver-address")

		client := client.NewClient(apiserverAddress, logger, _startOptions.Hostname)

		// Set up a signal channel to capture SIGTERM
		sigCh := make(chan os.Signal, 1)
		// interrupt signal sent from terminal
		signal.Notify(sigCh, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigCh, syscall.SIGTERM)
		signal.Notify(sigCh, syscall.SIGINT)

		// Start a goroutine to wait for the SIGTERM signal
		go func() {
			// Wait for the signal
			sig := <-sigCh
			logger.Sugar().Infof("received signal: %v", sig)

			// Perform cleanup or any necessary actions
			if err := client.Cleanup(context.TODO()); err != nil {
				logger.Sugar().Errorf("an error occurred in cleanup")
			}
			cancel() // Cancel the context
		}()

		if err := client.Start(ctx); err != nil {
			logger.Fatal(err.Error())
		}
	},
}

func init() {
	clientCmd.AddCommand(startCmd)
}
