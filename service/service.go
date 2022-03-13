package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SmartTriageFiap/notification/entities"
	"github.com/SmartTriageFiap/notification/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

var (
	once     = sync.Once{}
	instance Service
)

type Service interface {
	Notify(event entities.Event) error
}

type svc struct {
	db       repository.Repository
	mutex    sync.RWMutex
	partners []entities.Partner
	client   http.Client
}

func New(repo repository.Repository) Service {
	if instance == nil {
		once.Do(func() {
			svc := svc{
				db: repo,
				client: http.Client{
					Timeout: time.Second * 10,
				},
			}

			svc.gettingPartners()
			instance = &svc

		})
	}
	return instance
}

func (r svc) Notify(event entities.Event) error {
	for _, partner := range r.partners {
		b, err := json.Marshal(event)
		if err != nil {
			log.WithField("event", event).WithError(err).Warn("failed marshal event")
			continue
		}

		req, err := http.NewRequest(http.MethodPost, partner.Url, bytes.NewReader(b))
		if err != nil {
			log.WithField("event", event).WithError(err).Warn("failed creating request")
			continue
		}

		resp, err := r.client.Do(req)
		if err != nil {
			log.WithField("event", event).WithError(err).Warn("failed calling partner url")
			continue
		}

		log.Info(fmt.Sprintf("event sent successfully url: %s status: %d", partner.Url, resp.Status))
	}

	return nil
}

func (r svc) gettingPartners() {
	go func() {
		for {
			log.Info("getting partners")
			partners, err := r.db.GetPartnersAddresses()
			if err != nil {
				log.WithError(err).Warn("failed getting partners")
			}

			r.mutex.RLock()
			r.partners = partners
			r.mutex.RUnlock()
			<-time.After(time.Minute * 1)
		}
	}()
}
