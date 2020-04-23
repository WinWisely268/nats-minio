package main

// Import Go and NATS packages
import (
	"context"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"runtime"

	"github.com/nats-io/stan.go"
	"../pkg/event"
)

func main() {
	l := log.New(os.Stdout, "test-stan-minio", log.Lshortfile)
	ctx := context.Background()
	eventConf := &event.Conf{
		ClusterID: "test-cluster",
		Log: l,
		ID: "test-stan-publisher",
		ClientID: "test-stan-subscriber",
	}
	eventStreaming := event.New(eventConf)
	natsOpts := []nats.Option{nats.Name("Test NATS MINIO")}
	nc, err := nats.Connect("nats://localhost:4222", natsOpts...)
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
	runtime.Goexit()
}
