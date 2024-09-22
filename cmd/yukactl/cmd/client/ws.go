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

type wsOptions struct {
	Hostname string `flag:"ws-hostname"`
}

var _wsOptions wsOptions

// wsCmd represents the start command
var wsCommand = &cobra.Command{
	Use:   "ws",
	Short: "Starts up a websocket connection",
	Long: `Starts a websocket connection to the server using default configuration. 
Run "yukactl client ws --help" for more information.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateAndUnmarshal(cmd, &_wsOptions, validationFns); err != nil {
			log.Fatalln(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		logger, err := utils.GetLogger()
		if err != nil {
			logger.Fatal(err.Error())
		}

		client := client.NewWsWrapper(logger)

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
			logger.Sugar().Infof("Received signal: %v", sig)

			// Perform cleanup or any necessary actions
			if err := client.Cleanup(context.TODO()); err != nil {
				logger.Sugar().Errorf("An error occurred in cleanup", err)
			}
			cancel() // Cancel the context
		}()

		if err := client.StartWs(ctx, _wsOptions.Hostname); err != nil {
			logger.Fatal(err.Error())
		}
	},
}

func init() {
	wsCommand.PersistentFlags().StringP("ws-hostname", "w", "localhost:8082", "Hostname of the websocket server")
	clientCmd.AddCommand(wsCommand)
}
