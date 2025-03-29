package postgres

import (
	"errors"
	"log/slog"
)

const (
	postTable    = "posts"
	commentTable = "comments"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig(host, port, user, password, dbName, sslMode string) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

func (c *Config) Validate() error {
	if c.Host == "" {
		slog.Error("db: host is empty")
		return errors.New("db: host is empty")
	}

	if c.Port == "" {
		slog.Error("db: port is empty")
		return errors.New("db: port is empty")
	}

	if c.User == "" {
		slog.Error("db: user is empty")
		return errors.New("db: user is empty")
	}

	if c.Password == "" {
		slog.Error("db: password is empty")
		return errors.New("db: password is empty")
	}

	if c.DBName == "" {
		slog.Error("db: dbName is empty")
		return errors.New("db: dbName is empty")
	}

	if c.SSLMode == "" {
		slog.Error("db: sslMode is empty")
		return errors.New("db: sslMode is empty")
	}

	return nil
}
