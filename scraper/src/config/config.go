package config

import "os"

// Config is representation of a OS Env values
type Config struct {
	BILETIX_URL    string
	PASSO_URL      string
	KULTURIST_URL  string
	AMPQ_URI       string
	SCHEDULE_ARRAY string
	GRPC_ADDR      string
	GRPC_ADDR_PORT string
	EXCHANGE_NAME  string
	QUEUE_NAME     string
	ROUTING_KEY    string
	CONSUMER_TAG   string
}

// GetConfig returns values based on given os.Getenv("ENV")
func GetConfig() Config {
	if os.Getenv("ENV") == "dev" {
		return Config{
			BILETIX_URL:    os.Getenv("BILETIX_URL_DEV"),
			PASSO_URL:      os.Getenv("PASSO_URL_DEV"),
			KULTURIST_URL:  os.Getenv("KULTURIST_URL_DEV"),
			AMPQ_URI:       os.Getenv("AMPQ_URI_DEV"),
			SCHEDULE_ARRAY: os.Getenv("SCHEDULE_ARRAY_DEV"),
			GRPC_ADDR:      os.Getenv("GRPC_ADDR_DEV"),
			GRPC_ADDR_PORT: os.Getenv("GRPC_ADDR_PORT_DEV"),
			EXCHANGE_NAME:  os.Getenv("EXCHANGE_NAME_DEV"),
			QUEUE_NAME:     os.Getenv("QUEUE_NAME_DEV"),
			ROUTING_KEY:    os.Getenv("ROUTING_KEY_DEV"),
			CONSUMER_TAG:   os.Getenv("CONSUMER_TAG_DEV"),
		}
	} else if os.Getenv("ENV") == "test" {
		return Config{
			BILETIX_URL:    "https://www.biletix.com/solr/tr/select/",
			AMPQ_URI:       "amqp://guest:guest@localhost:5672",
			KULTURIST_URL:  "",
			SCHEDULE_ARRAY: `[{"provider":0,"timeType":0,"timeCount":1,"category":0,"datePeriod":0,"city":"İstanbul"},{"provider":1,"timeType":1,"timeCount":1,"category":0,"datePeriod":0,"city":"101"}]`,
			GRPC_ADDR:      "",
			GRPC_ADDR_PORT: "",
			EXCHANGE_NAME:  "events-test",
			QUEUE_NAME:     "delete-queue-test",
			ROUTING_KEY:    "routing-key-test",
			CONSUMER_TAG:   "consumer-test",
		}
	} else {
		return Config{
			BILETIX_URL:    os.Getenv("BILETIX_URL_PROD"),
			PASSO_URL:      os.Getenv("PASSO_URL_PROD"),
			KULTURIST_URL:  os.Getenv("KULTURIST_URL_PROD"),
			AMPQ_URI:       os.Getenv("AMPQ_URI_PROD"),
			SCHEDULE_ARRAY: os.Getenv("SCHEDULE_ARRAY_PROD"),
			GRPC_ADDR:      os.Getenv("GRPC_ADDR_PROD"),
			GRPC_ADDR_PORT: os.Getenv("GRPC_ADDR_PORT_PROD"),
			EXCHANGE_NAME:  os.Getenv("EXCHANGE_NAME_PROD"),
			QUEUE_NAME:     os.Getenv("QUEUE_NAME_PROD"),
			ROUTING_KEY:    os.Getenv("ROUTING_KEY_PROD"),
			CONSUMER_TAG:   os.Getenv("CONSUMER_TAG_PROD"),
		}
	}
}

// SetConfigForTesting sets the "ENV" value
func SetConfigForTesting() {
	os.Setenv("ENV", "test")
}
