package client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"yuka/internal/api/api_clients"
	"yuka/internal/client"
	"yuka/pkg/utils"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type startOptions struct {
	UserId         string   `flag:"user-id" validate:"required,uuid4"`
	OrganizationId string   `flag:"organization-id" validate:"required,uuid4"`
	TunnelIp       string   `flag:"tunnel-ip"`
	Hostname       string   `flag:"hostname"`
	ExposedPorts   []string `flag:"exposed-ports"`
}

var _startOptions startOptions

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts up the yuka daemon in 'client' mode",
	Long: `Starts up the yuka daemon in 'client' mode. 
This connects it to the yuka network for the given user id.

Run "yukad client start --help" for more information.`,
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
		transport := httptransport.New(apiserverAddress, "", nil)
		transport.DefaultAuthentication = httptransport.BasicAuth(os.Getenv("HTTP_USERNAME"), os.Getenv("HTTP_PASSWORD"))
		apiClient := api_clients.New(transport, strfmt.Default)

		client := client.NewClient(apiClient, logger, _startOptions.UserId, _startOptions.OrganizationId, _startOptions.TunnelIp, _startOptions.Hostname, _startOptions.ExposedPorts)

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
				logger.Sugar().Errorf("An error occurred in cleanup")
			}
			cancel() // Cancel the context
		}()

		if err := client.Start(ctx); err != nil {
			logger.Fatal(err.Error())
		}
	},
}

func init() {
	startCmd.Flags().StringP("user-id", "u", "", "Id of the user, currently just a UUID (E.g 'c5a24f48-3ebf-4e48-a351-22033b338c6b')")
	startCmd.Flags().StringP("organization-id", "o", "", "Id of the organization, currently just a UUID (E.g 'c5a24f48-3ebf-4e48-a351-22033b338c6b')")
	startCmd.Flags().StringP("tunnel-ip", "i", "", `Tunnel address of Wireguard (E.g 10.0.0.1:51280)`)
	startCmd.Flags().StringP("hostname", "n", "", "Hostname override")
	startCmd.Flags().StringSliceP("exposed-ports", "e", []string{}, "Exposed ports of known applications (E.g 8123)")

	clientCmd.AddCommand(startCmd)
}
