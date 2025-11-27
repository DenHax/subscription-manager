package postgres

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	URL string `env:"POSTGRES_URL" env-required:"true"`
}

type Storage struct {
	DB *sqlx.DB
}

func SetupConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("POSTGRES_URL is required and cannot be empty")
	}

	return &cfg, nil
}

func New(c Config) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", c.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println("Connection!")

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	fmt.Println("Close storage")
	return s.DB.Close()
}
