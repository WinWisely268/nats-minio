package main

// Import Go and NATS packages
import (
	"github.com/WinWisely268/nats-minio/pkg/event"
	"github.com/WinWisely268/nats-minio/pkg/svc"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "test-stan-minio", log.Lshortfile)
	eventConf := &event.Conf{
		ClusterID: "test-minio",
		Log:       l,
		ID:        "test-stan-publisher",
		ClientID:  "test-stan-subscriber",
	}
	eventStreaming := event.New(eventConf)
	natsOpts := []nats.Option{nats.Name("Test NATS MINIO")}
	nc, err := nats.Connect("nats://minio:minio@localhost:4222", natsOpts...)
	if err != nil {
		l.Fatalf("Unable to connect to NATS: %v", err)
	}
	defer nc.Close()
	if err = eventStreaming.Connect(stan.NatsConn(nc)); err != nil {
		l.Fatalf("Cannot connect to NATS Streaming server: %v", err)
	}
	subscriber := event.NewSubscriber(eventConf)
	err = subscriber.Connect(stan.NatsConn(nc))
	if err != nil {
		l.Fatalf("Cannnot subscribe to NATS streaming server: %v", err)
	}
	cfg := &svc.EventSvcCfg{
		Publisher:  eventStreaming,
		Subscriber: subscriber,
	}
	eventSvc := svc.New(cfg)
	eventSvc.Run()
}
