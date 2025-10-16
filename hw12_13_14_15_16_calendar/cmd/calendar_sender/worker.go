package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	configs "github.com/adettelle/hw/hw12_13_14_15_calendar/configs/sender"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	startCtx := context.Background()

	config, err := configs.New(&startCtx, "./configs/sender/sender_cfg.json")
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

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(logLevel),
	)

	fileForLog, err := os.OpenFile("/var/log/sender.log",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Println("Failed to create file")
		log.Fatal(err)
	}
	defer fileForLog.Close()

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(fileForLog),
		zap.NewAtomicLevelAt(logLevel),
	)

	var logg *zap.Logger
	// add DEBUG config param default false
	// logg = if DEBUG==true logg = logg := zap.New(zapcore.NewTee(core, fileCore))
	// else logg = zap.New(core)
	if config.Debug == "true" {
		logg = zap.New(zapcore.NewTee(core, fileCore))
	} else {
		logg = zap.New(core)
	}

	// logg := zap.New(zapcore.NewTee(core, fileCore)) //zap.New(core)
	logg.Info("LEVELS", zap.String("cfgLevel", config.Logger.Level), zap.String("actualLevel", logg.Level().String()))
	defer logg.Sync()

	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		logg.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logg.Fatal("Failed to open a channel", zap.Error(err))
	}
	defer ch.Close()

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

	msg, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logg.Fatal("Failed to register a consumer", zap.Error(err))
	}

	t, err := strconv.Atoi(config.WorkTicker)
	if err != nil {
		logg.Error("failed to parsecworkTicker", zap.Error(err))
	}
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for d := range msg {
		logg.Info("Received a message", zap.ByteString("message", d.Body))
	}
}
