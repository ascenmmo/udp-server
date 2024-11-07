package main

import (
	"context"
	"github.com/ascenmmo/udp-server/env"
	"github.com/ascenmmo/udp-server/pkg/start"
	"github.com/rs/zerolog"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	err := start.StartUDP(
		ctx,
		env.ServerAddress,
		env.TCPPort,
		env.UDPPort,
		env.TokenKey,
		env.MaxRequestPerSecond,
		5,
		logger,
		true)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start UDP server")
	}
}
