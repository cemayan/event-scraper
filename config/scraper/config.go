package scraper

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configer interface {
	LoadConfig(filename string) (*viper.Viper, error)
	ParseConfig(v *viper.Viper) (*AppConfig, error)
	GetConfig(configPath string) (*AppConfig, error)
}

type Config struct {
	viper *viper.Viper
}

// AppConfig is representation of a OS Env values
type AppConfig struct {
	SECRET         string
	GOARCH         string
	GOOS           string
	Postgresql     Postgresql
	Grpc           Grpc
	RabbitMQ       RabbitMQ
	Biletix        Biletix
	Kulturist      Kulturist
	Passo          Passo
	SCHEDULE_ARRAY string
}

type Postgresql struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	NAME     string
}

type Grpc struct {
	PORT string
	ADDR string
}

type RabbitMQ struct {
	AMPQ_URI      string
	QUEUE_NAME    string
	EXCHANGE_NAME string
	ROUTING_KEY   string
	CONSUMER_TAG  string
}

type Biletix struct {
	URL string
}

type Passo struct {
	URL string
}

type Kulturist struct {
	URL string
}

// Load config file from given path
func (cfg Config) LoadConfig(filename string) (*viper.Viper, error) {

	cfg.viper.SetConfigType("yaml")
	cfg.viper.AddConfigPath(".")
	cfg.viper.SetConfigName(filename)

	cfg.viper.AutomaticEnv()

	if err := cfg.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return cfg.viper, nil
}

// Parse config file
func (cfg Config) ParseConfig(v *viper.Viper) (*AppConfig, error) {
	var c AppConfig

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func (cfg Config) GetConfig(env string) (*AppConfig, error) {

	var path string
	if env == "dev" {
		path = "./config/scraper/config-dev"
	} else if env == "docker" {
		path = "./app/config/scraper/config-docker"
	} else if env == "prod" {
		path = "./app/config/scraper/config-prod"
	}

	v, err := cfg.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	_cfg, err := cfg.ParseConfig(v)
	if err != nil {
		return nil, err
	}
	return _cfg, nil
}

func NewConfig(viper *viper.Viper) Configer {
	return &Config{viper: viper}
}
