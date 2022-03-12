package main

import (
	"github.com/SmartTriageFiap/notification/config"
	"github.com/SmartTriageFiap/notification/listener"
	"github.com/SmartTriageFiap/notification/service"
	"github.com/arturmartini/envconfig"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	err := config.Configuration()
	if err != nil {
		log.WithError(err).Panic("error when try initialize config application")
	}
	listAddrs := strings.Split(envconfig.GetStr(config.KafkaAddr), " ")
	topic := envconfig.GetStr(config.KafkaTopic)
	svc := service.New()

	log.WithFields(log.Fields{
		"kafka-addr":  listAddrs,
		"kafka-topic": topic,
	}).Info("getting settings")

	log.WithError(listener.Start(svc, listAddrs, topic)).Panic("error when try start application")
}
