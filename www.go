package main

import (
	"context"
	"fmt"
	"net/http"

	"codeberg.org/reiver/go-field"

	"gamefed/cfg"
	"gamefed/srv/http"
	"gamefed/srv/log"

	_ "gamefed/www"
)

func www(ctx context.Context) <-chan struct{} {
	log := logsrv.Begin()
	defer log.End()

	done := make(chan struct{})

	var addr string = fmt.Sprintf(":%d", cfg.HTTPTcpPort())

	log.Trace(field.String("http-tcp-address", addr))

	server := &http.Server{
		Addr:    addr,
		Handler: &httpsrv.Mux,
	}

	go func() {
		defer close(done)
		log.Highlight(
			field.S("😈 www spawned"),
			field.String("http-tcp-address", server.Addr),
		)
		err := server.ListenAndServe()
		if nil != err && http.ErrServerClosed != err {
			log.Error(
				field.S("💀 www server error"),
				field.E(err),
			)
		}
		log.Highlight(
			field.S("😵 www died"),
			field.String("http-tcp-address", server.Addr),
		)
	}()

	go func() {
		<-ctx.Done()

		log.Highlight(
			field.S("👼 www killed"),
			field.String("http-tcp-address", server.Addr),
		)
		err := server.Shutdown(context.Background())
		if nil != err {
			log.Error(
				field.S("www server shutdown"),
				field.E(err),
			)
		}
	}()

	return done
}
