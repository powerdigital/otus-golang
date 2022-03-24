package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/logger"
	app_http "github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logger := logger.New(config.Logger)
	storage := storage.New(*config)
	app := app.New(logger, storage)
	server := app_http.NewServer(*logger, app)

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
