package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Control the miner service with RPC",
	Args: cobra.ExactValidArgs(1),
	ValidArgs: []string{"getStatus", "isPaused", "start", "stop"},
	Run: func(cmd *cobra.Command, args []string) {
		method := args[0]
		ctx := context.Background()
		clientpipe, err := rpc.DialIPC(ctx, NamedPipeName)
		defer clientpipe.Close()
		if err != nil {
			log.Fatal("Failed to connect to the named pipe")
		}
		switch method {
		case "getStatus":
			var result MiningStatus
			if err := clientpipe.Call(&result, "mining_getStatus"); err != nil {
				log.WithError(err).Fatal("Failed to call mining_getStatus")
			}
			fmt.Println("mining_getStatus:")
			fmt.Printf("\tisPaused:\t%t\n", result.IsPaused)
			fmt.Printf("\tisConnected:\t%t\n", result.IsConnected)
			fmt.Printf("\tpoolHostAndPort:\t%s\n", result.PoolHostAndPort)
			fmt.Printf("\tdurationConnected:\t%d\n", result.DurationConnected)
			fmt.Printf("\tblocksSubmitted:\t%d\n", result.BlocksSubmitted)
		case "isPaused":
			var result bool
			if err := clientpipe.Call(&result, "mining_isPaused"); err != nil {
				log.WithError(err).Fatal("Failed to call mining_isPaused")
			}
			fmt.Printf("mining_isPaused: %t\n", result)
		case "stop":
			clientpipe.Call(nil, "mining_stop")
		case "start":
			clientpipe.Call(nil, "mining_start")
		}
	},
}
