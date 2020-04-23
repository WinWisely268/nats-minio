package svc

import (
	"context"
	"github.com/nats-io/stan.go"
	"log"
	pb "github.com/WinWisely268/nats-minio/pkg/api"
)

type (
	EventSvcCfg struct {
		l *log.Logger
		Publisher *event.Publisher
		Subscriber *event.Subscriber
	}
	eventSvc struct {
		pub *event.Publisher
		sub *event.Subscriber
		l *log.Logger
	}
)

func New(conf *EventSvcCfg) *eventSvc {
	return &eventSvc{
		conf.Publisher,
		conf.Subscriber,
		conf.l,
	}
}

func (e *eventSvc) ListenNats(ctx context.Context) {
	mcb := func(msg *stan.Msg) {
		e.l.Printf("Received event: %s", string(msg.Data))
		var evt 
	}
}