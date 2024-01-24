package config

import (
	"auth-otp-go-grpc/pkg/utils"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	Database   `yaml:"database"`
	GrpcServer `yaml:"grpc_server"`
	JWT        `yaml:"jwt"`
	Smpp       `yaml:"smpp"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type GrpcServer struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	Idle_Timeout time.Duration `yaml:"idle_timeout"`
}

type Smpp struct {
	Smpp_Address          string `yaml:"smpp_address"`
	Smpp_User             string `yaml:"smpp_user"`
	Smpp_Password         string `yaml:"smpp_password"`
	Smpp_Src_Phone_Number string `yaml:"smpp_src_phone_number"`
}

type JWT struct {
	AccessSecretKey  string `yaml:"access_secret_key"`
	RefreshSecretKey string `yaml:"refresh_secret_key"`
}

func LoadConfig() Config {
	configPath := "./config/config.yaml"

	if configPath == "" {
		log.Fatalf("config path is not set or config file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %v", utils.Err(err))
	}

	return cfg
}
