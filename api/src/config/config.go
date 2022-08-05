package config

import "os"

// Config is representation of a OS Env values
type Config struct {
	DB_HOST       string
	DB_PORT       string
	DB_USER       string
	DB_PASSWORD   string
	DB_NAME       string
	AMPQ_URI      string
	SECRET        string
	GRPC_PORT     string
	EXCHANGE_NAME string
	QUEUE_NAME    string
	ROUTING_KEY   string
	CONSUMER_TAG  string
	AUTH_SERVER   string
}

// GetConfig returns values based on given os.Getenv("ENV")
func GetConfig() Config {
	if os.Getenv("ENV") == "dev" {
		return Config{
			DB_HOST:       os.Getenv("DB_HOST_DEV"),
			DB_PORT:       os.Getenv("DB_PORT_DEV"),
			DB_USER:       os.Getenv("DB_USER_DEV"),
			DB_PASSWORD:   os.Getenv("DB_PASSWORD_DEV"),
			DB_NAME:       os.Getenv("DB_NAME_DEV"),
			AMPQ_URI:      os.Getenv("AMPQ_URI_DEV"),
			SECRET:        os.Getenv("SECRET_DEV"),
			GRPC_PORT:     os.Getenv("GRPC_PORT_DEV"),
			EXCHANGE_NAME: os.Getenv("EXCHANGE_NAME_DEV"),
			QUEUE_NAME:    os.Getenv("QUEUE_NAME_DEV"),
			ROUTING_KEY:   os.Getenv("ROUTING_KEY_DEV"),
			CONSUMER_TAG:  os.Getenv("CONSUMER_TAG_DEV"),
			AUTH_SERVER:   os.Getenv("AUTH_SERVER_DEV"),
		}
	} else if os.Getenv("ENV") == "test" {
		return Config{
			DB_HOST:       "localhost",
			DB_PORT:       "5435",
			DB_USER:       "postgres",
			DB_PASSWORD:   "password",
			DB_NAME:       "scraper_db_test",
			AMPQ_URI:      "amqp://guest:guest@localhost:5672",
			SECRET:        "secret",
			GRPC_PORT:     "50051",
			EXCHANGE_NAME: "events-test",
			QUEUE_NAME:    "delete-queue-test",
			ROUTING_KEY:   "routing-key-test",
			CONSUMER_TAG:  "consumer-test",
			AUTH_SERVER:   "localhost:8109",
		}
	} else if os.Getenv("ENV") == "test_prod" {
		return Config{
			DB_HOST:       os.Getenv("DB_HOST_TEST_PROD"),
			DB_PORT:       os.Getenv("DB_PORT_TEST_PROD"),
			DB_USER:       os.Getenv("DB_USER_TEST_PROD"),
			DB_PASSWORD:   os.Getenv("DB_PASSWORD_TEST_PROD"),
			DB_NAME:       os.Getenv("DB_NAME_TEST_PROD"),
			AMPQ_URI:      os.Getenv("AMPQ_URI_TEST_PROD"),
			SECRET:        os.Getenv("SECRET_DEV_TEST_PROD"),
			GRPC_PORT:     os.Getenv("GRPC_PORT_TEST_PROD"),
			EXCHANGE_NAME: os.Getenv("EXCHANGE_TEST_PROD"),
			QUEUE_NAME:    os.Getenv("QUEUE_NAME_TEST_PROD"),
			ROUTING_KEY:   os.Getenv("ROUTING_KEY_TEST_PROD"),
			CONSUMER_TAG:  os.Getenv("CONSUMER_TAG_TEST_PROD"),
			AUTH_SERVER:   os.Getenv("AUTH_SERVER_TEST_PROD"),
		}
	} else if os.Getenv("ENV") == "prod" {
		return Config{
			DB_HOST:       os.Getenv("DB_HOST_PROD"),
			DB_PORT:       os.Getenv("DB_PORT_PROD"),
			DB_USER:       os.Getenv("DB_USER_PROD"),
			DB_PASSWORD:   os.Getenv("DB_PASSWORD_PROD"),
			DB_NAME:       os.Getenv("DB_NAME_PROD"),
			AMPQ_URI:      os.Getenv("AMPQ_URI_PROD"),
			SECRET:        os.Getenv("SECRET_PROD"),
			GRPC_PORT:     os.Getenv("GRPC_PORT_PROD"),
			EXCHANGE_NAME: os.Getenv("EXCHANGE_NAME_PROD"),
			QUEUE_NAME:    os.Getenv("QUEUE_NAME_PROD"),
			ROUTING_KEY:   os.Getenv("ROUTING_KEY_PROD"),
			CONSUMER_TAG:  os.Getenv("CONSUMER_TAG_PROD"),
			AUTH_SERVER:   os.Getenv("AUTH_SERVER_PROD"),
		}
	} else {
		return Config{}
	}
}
