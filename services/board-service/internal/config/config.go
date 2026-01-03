package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-default:"local"` // local, dev, prod
	Postgres PostgresConfig
	GRPC     GRPCConfig
	HTTP     HTTPConfig
}

type PostgresConfig struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
	// Можно разбить на Host, Port, User, если нужно, но URL удобнее для pgx
}

type GRPCConfig struct {
	Port    string        `env:"GRPC_PORT" env-default:":50051"`
	Timeout time.Duration `env:"GRPC_TIMEOUT" env-default:"5s"`
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT" env-default:":8080"`
}

func MustLoad() *Config {
	// Путь к конфиг-файлу. Можно брать из флага, но для простоты хардкодим или берем по умолчанию
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.env" // Дефолт для локальной разработки
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file %s not found, reading from ENV variables", configPath)
		// Если файла нет, cleanenv попытается прочитать просто ENV переменные системы
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		// Для продакшена, если файла нет, пробуем читать только ENV
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read config: %s", err)
		}
	}

	return &cfg
}
