package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs" //nolint:depguard
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/migrator"
	internalgrpc "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	config, err := configs.New(&startCtx)
	if err != nil {
		return err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	logLevel := zap.InfoLevel

	if config.Logger.Level != "" {
		logLevel, err = zapcore.ParseLevel(config.Logger.Level)
		if err != nil {
			log.Println("unable to set level")
			return err
		}
	}

	logg := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(logLevel),
	))
	logg.Info("LEVELS", zap.String("cfgLevel", config.Logger.Level), zap.String("actualLevel", logg.Level().String()))
	defer logg.Sync()

	logg.Info("Hello!")
	logg.Info(getVersion())

	storager, err := initStorager(config, logg)
	if err != nil {
		return err
	}

	calendar := app.New(logg, storager)

	var wg sync.WaitGroup

	server := internalhttp.NewServer(config, logg, calendar, storager)
	serverGRPC := internalgrpc.NewGRPCServer(config, logg, storager)

	go func() {
		s := <-ctx.Done()
		log.Printf("Got termination signal: %s. Graceful shutdown", s)

		stopCtx, cancel := context.WithTimeout(startCtx, time.Second*3)
		defer cancel()

		var err error
		if err = server.Stop(stopCtx); err != nil {
			logg.Error("failed to stop http server", zap.Error(err))
		}

		if err = serverGRPC.Close(); err != nil {
			logg.Error("failed to stop grpc server", zap.Error(err))
		}

		if err != nil {
			os.Exit(1)
		}

		<-stopCtx.Done()
		os.Exit(0)
	}()

	logg.Info("calendar is running...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(startCtx, logg); err != nil {
			logg.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(startCtx, logg); err != nil {
			logg.Fatal("failed to start grpc server", zap.Error(err))
			// return err
		}
	}()
	wg.Wait()
	return nil
}

// initStorager not only constructs, but also starts related processes
// depending on which storager we choose.
func initStorager(cfg *configs.Config, logg *zap.Logger) (app.Storager, error) {
	var storager app.Storager

	connStr := cfg.DBConnStr()

	if connStr != "" {
		db, err := database.Connect(connStr)
		if err != nil {
			return nil, err
		}

		storager = &sqlstorage.DBStorage{
			Ctx:  context.Background(),
			DB:   db,
			Logg: logg,
		}

		migrator.MustApplyMigrations(connStr, logg)
	} else {
		storager = memorystorage.New()
	}
	return storager, nil
}
