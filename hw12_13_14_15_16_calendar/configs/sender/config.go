package sender

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs"
)

const (
	defaultDBHost      = "localhost"
	defaultDBPort      = "9999"
	defaultDBUser      = "postgres"
	defaultDBPassword  = "123456"
	defaultDBName      = "calendar"
	defaultRabbitmqURL = "amqp://rmuser:rmpassword@localhost:5672/"
	defaultWokrTicker  = "5"
)

type Config struct {
	Logger     *LoggerConf `json:"logger"`
	Context    *context.Context
	Config     string // путь до json файла конфигурации
	DBHost     string `json:"dbhost"`
	DBPort     string `json:"dbport"`
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
	DBName     string `json:"dbname"`
	RabbitURL  string `json:"rabbiturl"`
	WorkTicker string `json:"workticker"`
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
		Config:     jsonPath,
		DBHost:     getEnvOrDefault("DBHOST", defaultDBHost),
		DBPort:     getEnvOrDefault("DBPORT", defaultDBPort),
		DBUser:     getEnvOrDefault("DBUSER", defaultDBUser),
		DBPassword: getEnvOrDefault("DBPASSWORD", defaultDBPassword),
		DBName:     getEnvOrDefault("DBNAME", defaultDBName),
		RabbitURL:  getEnvOrDefault("RURL", defaultRabbitmqURL),
		WorkTicker: getEnvOrDefault("WTICKER", defaultWokrTicker),
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
	if cfg.WorkTicker == "" {
		cfg.WorkTicker = defaultWokrTicker
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
	if cfg.WorkTicker == "" {
		cfg.WorkTicker = cfgFromJSON.WorkTicker
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
