package config

import (
	"github.com/arturmartini/envconfig"
)

const (
	KafkaAddr  = "KAFKA_ADDRS"
	KafkaTopic = "KAFKA_TOPIC"
)

func Configuration() error {
	config := &envconfig.Configuration{
		Envs: []string{
			KafkaAddr,
			KafkaTopic,
		},
		Required: []string{
			KafkaAddr,
			KafkaTopic,
		},
	}

	return envconfig.Initialize("", config)
}
