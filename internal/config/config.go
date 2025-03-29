package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

type Config struct {
	Server `yaml:"server"`
	DB     `yaml:"db"`
}

type Server struct {
	Port string `yaml:"port"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func initViper() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func InitConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading config file", err)
		return nil
	}

	err = initViper()
	if err != nil {
		slog.Error("Error initializing config", err)
		return nil
	}

	return &Config{
		Server: Server{
			Port: viper.GetString("server.port"),
		},

		DB: DB{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
	}
}
