package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	configs "github.com/adettelle/hw/hw12_13_14_15_calendar/configs/scheduler"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/migrator"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/pkg/database"
	_ "github.com/jackc/pgx/v5" // импортируем pgx для регистрации драйвера database/sql
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Notification struct {
	ID        string
	Title     string
	EventDate time.Time
	UserID    string
}

func main() {
	startCtx := context.Background()
	ctx, cancel := signal.NotifyContext(startCtx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, err := configs.New(&startCtx, "./configs/scheduler/scheduler_cfg.json")
	if err != nil {
		log.Fatal(err)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	logLevel := zap.InfoLevel

	if config.Logger.Level != "" {
		logLevel, err = zapcore.ParseLevel(config.Logger.Level)
		if err != nil {
			log.Println("unable to set level")
			log.Fatal(err)
		}
	}

	logg := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(logLevel),
	))
	logg.Info("LEVELS", zap.String("cfgLevel", config.Logger.Level), zap.String("actualLevel", logg.Level().String()))
	defer logg.Sync()

	// ----------------------------------
	// Создаем подключение к RabbitMQ
	conn, err := amqp.Dial(config.RabbitURL) // "amqp://rmuser:rmpassword@localhost:5672/"
	if err != nil {
		logg.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logg.Fatal("Failed to open a channel", zap.Error(err))
	}
	defer ch.Close()

	// объявляем queue для публикации сообщений
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logg.Fatal("Failed to declare a queue", zap.Error(err))
	}

	// ----------------------------------

	planner, err := initStorager(config, logg)
	if err != nil {
		logg.Fatal("failed to initialize storager", zap.Error(err))
	}

	t, err := strconv.Atoi(config.CollectTicker)
	if err != nil {
		logg.Error("failed to parsecworkTicker", zap.Error(err))
	}

	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctxPlanner, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			events, err := collectEvents(ctxPlanner, planner)
			if err != nil {
				logg.Error("failed to collect events", zap.Error(err))
			}
			ids, err := sendEvents(ctx, events, ch, q, logg)
			if err != nil {
				logg.Error("failed to send events", zap.Error(err))
			}
			_, err = planner.SetNotified(ctxPlanner, ids)
			if err != nil {
				logg.Error("failed to notify events", zap.Error(err))
			}
			err = planner.DeleteEvents(ctxPlanner)
			if err != nil {
				logg.Error("failed to delete events", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}

func collectEvents(ctx context.Context, planner Planner) ([]storage.EventToNotify, error) {
	events, err := planner.CollectEventsToNotify(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func sendEvents(ctxPlanner context.Context, events []storage.EventToNotify,
	ch *amqp.Channel, q amqp.Queue, logg *zap.Logger,
) ([]string, error) {
	var ids []string

	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}

		body := string(data)

		err = ch.PublishWithContext(ctxPlanner,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			})
		if err != nil {
			logg.Error("Failed to send an event", zap.Error(err))
		} else {
			ids = append(ids, event.ID)
		}
		log.Printf(" [x] Sent %s\n", body)
	}

	return ids, nil
}

func initStorager(cfg *configs.Config, logg *zap.Logger) (Planner, error) {
	var planner Planner

	connStr := cfg.DBConnStr()

	if connStr != "" {
		db, err := database.Connect(connStr)
		if err != nil {
			return nil, err
		}

		planner = &sqlstorage.DBStorage{
			Ctx:  context.Background(),
			DB:   db,
			Logg: logg,
		}

		migrator.MustApplyMigrations(connStr, logg)
	} else {
		planner = memorystorage.New()
	}
	return planner, nil
}

type Planner interface {
	SetNotified(ctx context.Context, ids []string) ([]string, error)
	CollectEventsToNotify(ctx context.Context) ([]storage.EventToNotify, error)
	DeleteEvents(ctx context.Context) error
}
