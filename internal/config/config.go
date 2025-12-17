package config

import (
	"fmt"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config — корневая структура
	Config struct {
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
		GRPC `yaml:"grpc"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
		User     string `env-required:"true" yaml:"user" env:"PG_USER"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		DBName   string `env-required:"true" yaml:"dbname" env:"PG_DBNAME"`
		SSLMode  string `yaml:"ssl_mode" env:"PG_SSL_MODE" env-default:"disable"`
		PoolMax  int    `yaml:"pool_max" env:"PG_POOL_MAX" env-default:"10"`
	}

	GRPC struct {
		UsersAddress  string `yaml:"users_address" env:"GRPC_USERS_ADDRESS" env-default:":50051"`
		BoardsAddress string `yaml:"boards_address" env:"GRPC_BOARDS_ADDRESS" env-default:":50052"`
	}
)

// NewConfig инициализирует конфигурацию
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// 1. Попытка прочитать из .env файла (удобно для локальной разработки)
	// Если файла нет — ошибка игнорируется, идем читать ENV системы
	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		// Если .env нет, читаем переменные окружения напрямую (для продакшена/Docker)
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	return cfg, nil
}

func (p *PG) ConnectionURL() string {
	// Формат: postgres://user:password@host:port/dbname?sslmode=disable
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%s", p.Host, p.Port),
		Path:   p.DBName,
	}

	q := u.Query()
	q.Set("sslmode", p.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
