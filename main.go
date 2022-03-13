package main

import (
	"github.com/SmartTriageFiap/notification/config"
	"github.com/SmartTriageFiap/notification/listener"
	"github.com/SmartTriageFiap/notification/repository"
	"github.com/SmartTriageFiap/notification/service"
	"github.com/arturmartini/envconfig"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	err := config.Initialize()
	if err != nil {
		log.WithError(err).Panic("error when try initialize config application")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	hosts := strings.Split(envconfig.GetStr(config.KafkaAddr), " ")
	topic := envconfig.GetStr(config.KafkaTopic)
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: hosts, Topic: topic})
	repo := repository.New()
	svc := service.New(repo)
	ltn := listener.New()

	log.WithError(ltn.Start(svc, reader, sigs)).Panic("error when try start application")
}
