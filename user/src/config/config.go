package config

import "os"

// Config is representation of a OS Env values
type Config struct {
	DB_HOST            string
	DB_PORT            string
	DB_USER            string
	DB_PASSWORD        string
	DB_NAME            string
	SECRET             string
	ERROR_AMQPS        string
	REDIS_ADDRESS      string
	REDIS_ADDRESS_PORT string
}

// GetConfig returns values based on given os.Getenv("ENV")
func GetConfig() Config {
	if os.Getenv("ENV") == "dev" {
		return Config{
			DB_HOST:            os.Getenv("DB_HOST_DEV"),
			DB_PORT:            os.Getenv("DB_PORT_DEV"),
			DB_USER:            os.Getenv("DB_USER_DEV"),
			DB_PASSWORD:        os.Getenv("DB_PASSWORD_DEV"),
			DB_NAME:            os.Getenv("DB_NAME_DEV"),
			REDIS_ADDRESS:      os.Getenv("REDIS_ADDRESS_DEV"),
			REDIS_ADDRESS_PORT: os.Getenv("REDIS_ADDRESS_PORT_DEV"),
			SECRET:             os.Getenv("SECRET_DEV"),
		}
	} else if os.Getenv("ENV") == "test" {
		return Config{
			DB_HOST:            "localhost",
			DB_PORT:            "5435",
			DB_USER:            "postgres",
			DB_PASSWORD:        "password",
			DB_NAME:            "scraper_db_test",
			REDIS_ADDRESS:      "localhost",
			REDIS_ADDRESS_PORT: "6379",
			SECRET:             "secret",
		}
	} else {
		return Config{
			DB_HOST:            os.Getenv("DB_HOST_PROD"),
			DB_PORT:            os.Getenv("DB_PORT_PROD"),
			DB_USER:            os.Getenv("DB_USER_PROD"),
			DB_PASSWORD:        os.Getenv("DB_PASSWORD_PROD"),
			DB_NAME:            os.Getenv("DB_NAME_PROD"),
			REDIS_ADDRESS:      os.Getenv("REDIS_ADDRESS_PROD"),
			REDIS_ADDRESS_PORT: os.Getenv("REDIS_ADDRESS_PORT_PROD"),
			SECRET:             os.Getenv("SECRET_PROD"),
		}
	}
}
