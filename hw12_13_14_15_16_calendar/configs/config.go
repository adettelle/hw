package configs

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

// флаг только -config
// если он есть, то конфиг файл + default
// если нет такого флага, то env + default

const (
	defaultAddress     = "localhost:8080"
	defaultGRPCAddress = "localhost:8082"
	defaultDBHost      = "localhost"
	defaultDBPort      = "9999"
	defaultDBUser      = "postgres"
	defaultDBPassword  = "123456"
	defaultDBName      = "calendar"
)

// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger      *LoggerConf `json:"logger"`
	Context     *context.Context
	Config      string // путь до json файла конфигурации по умолчанию /configs/cfg.json
	Address     string `json:"address"`
	GRPCAddress string `json:"grpcaddr"`
	DBHost      string `json:"dbhost"`
	DBPort      string `json:"dbport"`
	DBUser      string `json:"dbuser"`
	DBPassword  string `json:"dbpassword"`
	DBName      string `json:"dbname"`
}

type LoggerConf struct {
	Level string `json:"level"`
}

func New(ctx *context.Context, ignoreFlags bool, jsonPath string) (*Config, error) {
	var cfg *Config

	flagConfig := flag.String("config", "", "path to file with config parametrs")
	flag.Parse()

	if *flagConfig != "" {
		cfgFromJSON, err := ReadCfgJSON[Config](*flagConfig)
		if err != nil {
			return nil, err
		}
		cfgFromJSON.applyDefauls()
		cfgFromJSON.Context = ctx

		ensureAddrFLagIsCorrect(cfgFromJSON.Address)
		ensureAddrFLagIsCorrect(cfgFromJSON.GRPCAddress)
		ensureHostFlagIsCorrect(*cfgFromJSON.Context, cfgFromJSON.DBHost)
		ensurePortFlagIsCorrect(cfgFromJSON.DBPort)

		return cfgFromJSON, nil
	} else {
		cfg = &Config{
			Logger: &LoggerConf{
				Level: "INFO",
			},
			Address:     getEnvOrDefault("DBHOST", defaultAddress),
			GRPCAddress: getEnvOrDefault("DBHOST", defaultGRPCAddress),
			Config:      "", // *flagConfig,

			DBHost:     getEnvOrDefault("DBHOST", defaultDBHost),
			DBPort:     getEnvOrDefault("DBPORT", defaultDBPort),
			DBUser:     getEnvOrDefault("DBUSER", defaultDBUser),
			DBPassword: getEnvOrDefault("DBPASSWORD", defaultDBPassword),
			DBName:     getEnvOrDefault("DBNAME", defaultDBName),
		}
	}

	cfg.Context = ctx

	ensureAddrFLagIsCorrect(cfg.Address)
	ensureAddrFLagIsCorrect(cfg.GRPCAddress)
	ensureHostFlagIsCorrect(*cfg.Context, cfg.DBHost)
	ensurePortFlagIsCorrect(cfg.DBPort)

	return cfg, nil
}

// заполняем структуру конфига из default.
func (cfg *Config) applyDefauls() {
	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}
	if cfg.GRPCAddress == "" {
		cfg.GRPCAddress = defaultGRPCAddress
	}

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
}

func getEnvOrDefault(envName string, defaultVal string) string {
	res := os.Getenv(envName)

	if res != "" {
		return res
	}
	return defaultVal
}

func ensureAddrFLagIsCorrect(addr string) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}
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
