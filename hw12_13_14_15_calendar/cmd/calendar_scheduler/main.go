package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/client/kafka"                 //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	memorystorage "github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yml", "Path to configuration file")
}

func main() {
	// docker run -p 9092:9092 apache/kafka-native:3.9.0
	flag.Parse()

	if flag.Arg(0) == "version" {
		config.PrintVersion()
		return
	}

	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.File)
	defer logg.Close()

	var storage app.Storage
	if config.Storage == "sql" {
		logg.Info("create sql storage, connecting to server...")
		dbStorage := sqlstorage.New(config.DB.Driver, config.DB.Dsn)
		err := dbStorage.Connect(context.Background())
		if err != nil {
			logg.Error("failed to connect to db: " + err.Error())
			os.Exit(1) //nolint:gocritic
		}
		storage = dbStorage
	} else {
		logg.Info("create memory storage")
		storage = memorystorage.New()
	}
	app.New(logg, storage)

	kafka := kafka.New([]string{net.JoinHostPort(config.Kafka.Host, strconv.Itoa(config.Kafka.Port))})

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		dbStorage, ok := storage.(*sqlstorage.Storage)
		if ok {
			if err := dbStorage.Close(ctx); err != nil {
				logg.Error("failed to close database: " + err.Error())
			}
		}
	}()

	logg.Info("calendar scheduler is starting...")
	logg.Info(config.Kafka.Host + ":" + strconv.Itoa(config.Kafka.Port))
	logg.Info(config.Kafka.Topic)

	if err := kafka.Connect(ctx); err != nil {
		logg.Error("failed to connect to kafka: " + err.Error())
		return
	}
}
