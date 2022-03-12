package listener

import (
	"context"
	"fmt"
	"github.com/SmartTriageFiap/notification/service"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

func Start(svc service.Service, hosts []string, topic string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: hosts,
		Topic:   topic,
	})

	ctx := context.Background()
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Warn(fmt.Sprintf("failed to read message partition: %s offset: %s", m.Partition, m.Offset), err)
		}

		if err := svc.Notify(string(m.Value)); err != nil {
			log.Warn(fmt.Sprintf("failed to notify partition: %s offset: %s", m.Partition, m.Offset), err)
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Warn(fmt.Sprintf("failed to commit messages partition: %s offset: %s", m.Partition, m.Offset), err)
		}
	}

}
