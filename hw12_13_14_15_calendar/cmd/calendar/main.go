package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	internalhttp "github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/server/http"     //nolint:depguard
	memorystorage "github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.File)
	defer logg.Close()

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
