package service

import "sync"

var (
	once     = sync.Once{}
	instance Service
)

type Service interface {
	Notify(string) error
}

type svc struct{}

func New() Service {
	if instance == nil {
		once.Do(func() {
			instance = &svc{}
		})
	}
	return instance
}

func (r svc) Notify(string) error {
	return nil
}
