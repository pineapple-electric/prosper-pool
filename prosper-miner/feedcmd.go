package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:    "feed",
	Short:  "Show notifications from the miner service",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		clientpipe, err := rpc.DialIPC(ctx, NamedPipeName)
		if err != nil {
			log.Fatal("Failed to connect to the named pipe")
		}
		defer clientpipe.Close()
		hashRateChannel := make(chan float64)
		submissionChannel := make(chan int)

		subRate, err := clientpipe.Subscribe(ctx, "mining", hashRateChannel, "hashRateSubscription")
		if err != nil {
			log.WithError(err).Fatal("Failed to subscribe to mining:hashRateSubscription")
		}
		subSubmission, err := clientpipe.Subscribe(ctx, "mining", submissionChannel, "submissionSubscription")
		if err != nil {
			log.WithError(err).Fatal("Failed to subscribe to mining:submissionSubscription")
		}

		go func() {
			for {
				select {
				case hashRate := <-hashRateChannel:
					fmt.Printf("Current hash rate: %f\n", hashRate)
				case <-submissionChannel:
					fmt.Println("A share was submitted")
				}
			}
		}()

		// Allow the user to exit by pressing Enter
		fmt.Println("Press Enter to terminate the feed")
		keyboardReader := bufio.NewReader(os.Stdin)
		keyboardReader.ReadString('\n')
		subRate.Unsubscribe()
		subSubmission.Unsubscribe()
		return
	},
}
