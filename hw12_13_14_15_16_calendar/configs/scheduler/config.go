package scheduler

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs"
)

const (
	defaultDBHost        = "localhost"
	defaultDBPort        = "9999"
	defaultDBUser        = "postgres"
	defaultDBPassword    = "123456"
	defaultDBName        = "test_db" // calendar
	defaultRabbitmqURL   = "amqp://rmuser:rmpassword@localhost:5672/"
	defaultCollectTicker = "5"
)

type Config struct {
	Logger        *LoggerConf `json:"logger"`
	Context       *context.Context
	Config        string // путь до json файла конфигурации
	DBHost        string `json:"dbhost"`
	DBPort        string `json:"dbport"`
	DBUser        string `json:"dbuser"`
	DBPassword    string `json:"dbpassword"`
	DBName        string `json:"dbname"`
	RabbitURL     string `json:"rabbiturl"`
	CollectTicker string `json:"collectticker"`
}

type LoggerConf struct {
	Level string `json:"level"`
}

// далее проверяем, если есть json файл и заполняем структкуру конфига оттуда;
// далее проверяем, если поле не заполнено, заполняем по default.
func New(ctx *context.Context, jsonPath string) (*Config, error) {
	cfg := Config{
		Logger: &LoggerConf{
			Level: "INFO",
		},
		Config:        jsonPath,
		DBHost:        getEnvOrDefault("DBHOST", defaultDBHost),
		DBPort:        getEnvOrDefault("DBPORT", defaultDBPort),
		DBUser:        getEnvOrDefault("DBUSER", defaultDBUser),
		DBPassword:    getEnvOrDefault("DBPASSWORD", defaultDBPassword),
		DBName:        getEnvOrDefault("DBNAME", defaultDBName),
		RabbitURL:     getEnvOrDefault("RURL", defaultRabbitmqURL),
		CollectTicker: getEnvOrDefault("TICKER", defaultCollectTicker),
	}
	cfgFromJSON, err := configs.ReadCfgJSON[Config](jsonPath)
	if err != nil {
		return nil, err
	}
	cfg.applyConfigFromJSON(cfgFromJSON)

	cfg.Context = ctx

	cfg.applyDefauls()

	ensureHostFlagIsCorrect(*cfg.Context, cfg.DBHost)
	ensurePortFlagIsCorrect(cfg.DBPort)

	return &cfg, nil
}

// заполняем структуру конфига из default.
func (cfg *Config) applyDefauls() {
	if cfg.DBHost == "" {
		cfg.DBHost = defaultDBHost
	}

	if cfg.DBPort == "" {
		cfg.DBPort = defaultDBPort
	}
	if cfg.DBUser == "" {
		cfg.DBUser = defaultDBUser
	}
	if cfg.DBName == "" {
		cfg.DBName = defaultDBName
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = defaultDBPassword
	}
	if cfg.RabbitURL == "" {
		cfg.RabbitURL = defaultRabbitmqURL
	}
	if cfg.CollectTicker == "" {
		cfg.CollectTicker = defaultCollectTicker
	}
}

// проверяем, если есть json файл и дополняем структкуру конфига оттуда.
func (cfg *Config) applyConfigFromJSON(cfgFromJSON *Config) {
	if cfg.DBHost == "" {
		cfg.DBHost = cfgFromJSON.DBHost
	}
	if cfg.DBPort == "" {
		cfg.DBPort = cfgFromJSON.DBPort
	}
	if cfg.DBUser == "" {
		cfg.DBUser = cfgFromJSON.DBUser
	}
	if cfg.DBName == "" {
		cfg.DBName = cfgFromJSON.DBName
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = cfgFromJSON.DBPassword
	}
	if cfg.Logger == nil {
		cfg.Logger = cfgFromJSON.Logger
	}
	if cfg.RabbitURL == "" {
		cfg.RabbitURL = cfgFromJSON.RabbitURL
	}
	if cfg.CollectTicker == "" {
		cfg.CollectTicker = cfgFromJSON.CollectTicker
	}
}

func getEnvOrDefault(envName string, defaultVal string) string {
	res := os.Getenv(envName)
	if res != "" {
		return res
	}
	return defaultVal
}

func ensureHostFlagIsCorrect(ctx context.Context, host string) {
	resolver := net.Resolver{}

	addrs, err := resolver.LookupHost(ctx, host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("host in ensureHostFlagIsCorrect:", addrs)
}

func ensurePortFlagIsCorrect(port string) {
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("xxx invalid port:", err)
	}
}

// defaultDBParams = "host=localhost port=9999 user=postgres password=123456 dbname=calendar sslmode=disable"
// DBConnStr constructs and returns the PostgreSQL database connection string.
func (cfg *Config) DBConnStr() string {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	return dbParams
}
