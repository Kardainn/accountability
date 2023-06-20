package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Kardainn/accountability/backend/config"
	"github.com/Kardainn/accountability/backend/database"
	"github.com/Kardainn/accountability/backend/server"
)

func main() {
	// create empty context
	ctx := context.Background()
	// Receiving a new context
	ctx, err := config.Init(ctx)
	if err != nil {
		panic(err)
	}
	server := server.Create(ctx)
	httpChan := make(chan error)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		httpChan <- fmt.Errorf("received Interrupt signal")
	}()
	server.Start(httpChan)
	database.ConnectDB(ctx)
	err = <-httpChan
	println(err.Error())
}
