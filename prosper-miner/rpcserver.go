package main

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

var NamedPipeName = `\\.\pipe\ProsperPoolService`

type MinerRPCService struct {
	m *Mining
}

// RPC methods

// Exposed as mining_getStatus
func (s MinerRPCService) GetStatus() *MiningStatus {
	var status *MiningStatus
	status = s.m.GetStatus()
	return status
}

// Exposed as mining_isPaused
func (s MinerRPCService) IsPaused() bool {
	return s.m.IsPaused()
}

// Exposed as mining_start
func (s MinerRPCService) Start() {
	s.m.Start()
}

// Exposed as mining_stop
func (s MinerRPCService) Stop() {
	s.m.Stop()
}

// Subscriptions

func (s MinerRPCService) HashRateSubscription(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}
	subscription := notifier.CreateSubscription()
	go func() {
		if hashRateChannel := s.m.GetHashRateChannel(); hashRateChannel != nil {
			for {
				select {
				case <-notifier.Closed():
					// Client closed the connection
				case <-subscription.Err():
					// The client unsubscribed
				case i := <-hashRateChannel:
					notifier.Notify(subscription.ID, i)
				}
			}
		}
	}()

	return subscription, nil
}

func (s MinerRPCService) SubmissionSubscription(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}
	subscription := notifier.CreateSubscription()
	go func() {
		if submissionChannel := s.m.GetSubmissionChannel(); submissionChannel != nil {
			for {
				select {
				case <-notifier.Closed():
					// Client closed the connection
				case <-subscription.Err():
					// The client unsubscribed
				case i := <-submissionChannel:
					notifier.Notify(subscription.ID, i)
				}
			}
		}
	}()

	return subscription, nil
}

func startRPCServer(mining *Mining) {
	newApi := rpc.API{}
	newApi.Namespace = "mining"
	newApi.Version = "1"
	newApi.Service = MinerRPCService{mining}
	newApi.Public = true
	_, _, err := rpc.StartIPCEndpoint(NamedPipeName, []rpc.API{newApi})
	if err != nil {
		// This may happen if the service is already running
		log.WithError(err).Fatal("Failed to start the RPC Server on the named pipe")
	}
}
