package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"in-mem-io-task-manager/internal/application"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	app, err := application.New(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go func() {
		<-signalCh
		if err := app.Close(ctx); err != nil {
			log.Println(err)
		}

		cancel()
	}()

	if err := app.Start(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	<-ctx.Done()
}
