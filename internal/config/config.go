package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Host       string
	Port       string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func Create() (*Config, error) {
	cmd, _ := os.Getwd()
	err := godotenv.Load(cmd + "\\.env")
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:       os.Getenv("HOST"),
		Port:       os.Getenv("PORT"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
	}, nil
}

func (c *Config) getHost() string {
	return c.Host
}

func (c *Config) getPort() string {
	return c.Port
}
