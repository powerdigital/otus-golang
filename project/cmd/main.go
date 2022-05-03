package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/powerdigital/project/internal/app"
	internalhttp "github.com/powerdigital/project/internal/server/http"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	app := app.NewApp()
	server := internalhttp.NewServer(app, *config)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			return
		}
	}()

	fmt.Println("app is running...")

	if err := server.Start(ctx); err != nil {
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
