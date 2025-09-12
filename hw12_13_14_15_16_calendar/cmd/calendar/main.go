package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs" //nolint:depguard
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/database"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/migrator"
	internalhttp "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/memory"
	"go.uber.org/zap"
)

func main() {
	err := initialize()
	if err != nil {
		log.Fatal(err)
	}
}

func initialize() error {
	startCtx := context.Background()
	ctx, cancel := signal.NotifyContext(startCtx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, err := configs.New(&startCtx, true, "./configs/config.yaml")
	if err != nil {
		return err
	}

	logg, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logg.Sync()

	logg.Info("Hello!")
	logg.Info(getVersion())

	connStr := config.DBConnStr()

	migrator.MustApplyMigrations(connStr)

	db, err := database.Connect(connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config, logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(startCtx, time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server", zap.Error(err))
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(startCtx); err != nil { // ctx ????
		logg.Error("failed to start http server", zap.Error(err))
		return err
	}
	return nil
}
