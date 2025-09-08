package configs

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/helpers"
)

const (
	defaultAddress    = "localhost:8080"
	defaultDBHost     = "localhost"
	defaultDBPort     = "9999"
	defaultDBUser     = "postgres"
	defaultDBPassword = "123456"
	defaultDBName     = "calendar"
	// defaultDBParams = "host=localhost port=9999 user=postgres password=password dbname=calendar sslmode=disable"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf
	Config     string // путь до json файла конфигурации по умолчанию /config/cfg.json
	Address    string `json:"address"`
	DBHost     string `json:"dbhost"`
	DBPort     string `json:"dbport"`
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
	DBName     string `json:"dbname"`
	//DBParams   string `json:"database_dsn"` // TODO HELP
}

type LoggerConf struct {
	Level string
	// TODO
}

func InitFlags() *Config {
	flagAddr := flag.String("a", "", "Net address localhost:port")
	flagConfig := flag.String("config", "", "path to file with config parametrs")

	flagDBHost := flag.String("h", "", "dbhost")
	flagDBPort := flag.String("p", "", "dbport")
	flagDBUser := flag.String("u", "", "db user")
	flagDBPassword := flag.String("password", "", "db password")
	flagDBName := flag.String("n", "", "db name")
	// flagDBParams := flag.String("d", "", "db connection params")

	flag.Parse()

	cfg := Config{
		Logger:  LoggerConf{}, // TODO
		Address: getAddr(flagAddr),
		Config:  getConfig(flagConfig),

		DBHost:     getHost(flagDBHost),
		DBPort:     getPort(flagDBPort),
		DBUser:     getUser(flagDBUser),
		DBPassword: getPassword(flagDBPassword),
		DBName:     getName(flagDBName),
		//DBParams:   getDBParams(flagDBParams),
	}
	return &cfg
}

// приоритет:
// сначала проверяем флаги и заполняем структуру конфига оттуда;
// потом проверяем переменные окружения и перезаписываем структуру конфига оттуда;
// далее проверяем, если есть json файл и дополняем структкуру конфига оттуда.
func New(ignoreFlags bool, jsonPath string) (*Config, error) {
	var cfg *Config

	if !ignoreFlags {
		cfg = InitFlags()
	} else {
		cfg = &Config{Config: jsonPath}
	}
	if cfg.Config != "" {
		cfgFromJSON, err := helpers.ReadCfgJSON[Config](cfg.Config)
		if err != nil {
			return nil, err
		}
		if cfg.Address == "" {
			cfg.Address = cfgFromJSON.Address
		}
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
	}

	if cfg.Address == "" {
		cfg.Address = defaultAddress
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

	// if cfg.DBParams == "" {
	// 	cfg.DBParams = defaultDBParams
	// }
	ensureAddrFLagIsCorrect(cfg.Address)
	ensureHostFlagIsCorrect(cfg.DBHost)
	ensurePortFlagIsCorrect(cfg.DBPort)

	return cfg, nil
}

func getConfig(flagConfig *string) string {
	config := os.Getenv("CONFIG")
	if config != "" {
		return config
	}
	return *flagConfig
}

func getAddr(flagAddr *string) string {
	addr := os.Getenv("ADDRESS")
	if addr != "" {
		return addr
	}
	return *flagAddr
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

// func getDBParams(flagDBParams *string) string {
// 	envDBParams := os.Getenv("DATABASE_DSN")

// 	if envDBParams != "" {
// 		return envDBParams
// 	}

// 	return *flagDBParams
// }

func getHost(flagHost *string) string {
	host := os.Getenv("DBHOST")
	if host != "" {
		return host
	}
	return *flagHost
}

func ensureHostFlagIsCorrect(host string) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("host in ensureHostFlagIsCorrect:", addrs)
}

func getPort(flagPort *string) string {
	port := os.Getenv("DBPORT")
	if port != "" {
		return port
	}
	return *flagPort
}

func ensurePortFlagIsCorrect(port string) {
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("invalid port:", err)
	}
}

func getUser(flagDBUser *string) string {
	user := os.Getenv("DBUSER")
	if user != "" {
		return user
	}
	return *flagDBUser
}
func getPassword(flagDBPassword *string) string {
	pwd := os.Getenv("DBPASSWORD")
	if pwd != "" {
		return pwd
	}
	return *flagDBPassword
}
func getName(flagDBName *string) string {
	name := os.Getenv("DBNAME")
	if name != "" {
		return name
	}
	return *flagDBName
}

// defaultDBParams = "host=localhost port=9999 user=postgres password=123456 dbname=calendar sslmode=disable"
// DBConnStr constructs and returns the PostgreSQL database connection string.
func (cfg *Config) DBConnStr() string {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	return dbParams
}
