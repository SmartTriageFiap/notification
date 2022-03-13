package config

import (
	"github.com/arturmartini/envconfig"
)

const (
	KafkaAddr  = "KAFKA_ADDRS"
	KafkaTopic = "KAFKA_TOPIC"
	MongoUser  = "MONGO_USER"
	MongoPass  = "MONGO_PASS"
	MongoAddr  = "MONGO_ADDR"
	MongoPort  = "MONGO_PORT"
)

func Initialize() error {
	config := &envconfig.Configuration{
		Envs: []string{
			KafkaAddr,
			KafkaTopic,
			MongoUser,
			MongoPass,
			MongoAddr,
			MongoPort,
		},
		Required: []string{
			KafkaAddr,
			KafkaTopic,
			MongoUser,
			MongoPass,
			MongoAddr,
			MongoPort,
		},
	}

	return envconfig.Initialize("", config)
}
