package main

import (
	"context"
	"os/signal"
	"syscall"

	"gamefed/srv/log"
)

func run() {
	log := logsrv.Begin()
	defer log.End()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(sigCtx)
	defer cancel()

	wwwDone := www(ctx)

	select {
	case <-wwwDone:
	}

	cancel()

	<-wwwDone
}
