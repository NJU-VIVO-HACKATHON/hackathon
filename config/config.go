package config

import (
	"time"
)

type Database struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db-name"`
}

type Config struct {
	Database Database     `yaml:"database"`
	Jwt      JwtConfig    `yaml:"jwt"`
	Server   ServerConfig `yaml:"server"`
}

type JwtConfig struct {
	SecretKey       string        `yaml:"secret-key"`
	ExpiresDuration time.Duration `yaml:"expires-duration"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
