package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address      string        `yaml:"address" env-default:"localhost:8080" yaml-required:"true"`
	SSLMode      string        `yaml:"ssl_mode" env-default:"disable"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func SetupConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &cfg, nil
}

type Server struct {
	httpServer *http.Server
}

func New(scfg Config, handler http.Handler) *Server {
	srv := new(Server)
	srv.httpServer = &http.Server{
		Addr:           scfg.Address,
		Handler:        handler,
		ReadTimeout:    scfg.ReadTimeout,
		WriteTimeout:   scfg.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	return srv
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
