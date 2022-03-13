package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SmartTriageFiap/notification/entities"
	"github.com/SmartTriageFiap/notification/service"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	onde = sync.Once{}
	ltn  Listener
)

type Listener interface {
	Start(service.Service, *kafka.Reader, chan os.Signal) error
}

type listener struct {
	run               bool
	inProcessing      bool
	mutexRun          sync.Mutex
	mutexInProcessing sync.Mutex
}

func New() Listener {
	if ltn == nil {
		onde.Do(func() {
			ltn = &listener{
				run:               true,
				inProcessing:      false,
				mutexRun:          sync.Mutex{},
				mutexInProcessing: sync.Mutex{},
			}
		})
	}
	return ltn
}

func (r listener) Start(svc service.Service, reader *kafka.Reader, stop chan os.Signal) error {
	if reader == nil {
		log.Panic("failed reader is null")
	}

	go r.gracefullShutdown(stop)

	ctx := context.Background()
	for r.run {
		r.mutexInProcessing.Lock()
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Warn(fmt.Sprintf("failed to read message partition: %s offset: %s", m.Partition, m.Offset), err)
		}
		r.inProcessing = true

		event, err := translateByteToEvent(m.Value)
		if err != nil {
			log.Warn(fmt.Sprintf("failed to translate message partition: %s offset: %s", m.Partition, m.Offset), err)
		}

		if err := svc.Notify(event); err != nil {
			log.Warn(fmt.Sprintf("failed to notify partition: %s offset: %s", m.Partition, m.Offset), err)
		}

		if err := reader.CommitMessages(ctx, m); err != nil {
			log.Warn(fmt.Sprintf("failed to commit messages partition: %s offset: %s", m.Partition, m.Offset), err)
		}
		r.inProcessing = false
		r.mutexInProcessing.Unlock()
	}

	return nil
}

func translateByteToEvent(value []byte) (entities.Event, error) {
	event := entities.Event{}
	err := json.Unmarshal(value, &event)
	return event, err
}

func (r listener) gracefullShutdown(stop chan os.Signal) {
	sig := <-stop
	log.WithField("signal", sig).Info("start graceful shutdown")
	r.mutexRun.Lock()
	r.run = false
	r.mutexRun.Unlock()
	for r.inProcessing {
		<-time.After(time.Second * 5)
	}
	os.Exit(0)
}
