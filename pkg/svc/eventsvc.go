package svc

import (
	"context"
	"errors"
	pb "github.com/WinWisely268/nats-minio/pkg/api"
	"github.com/WinWisely268/nats-minio/pkg/event"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nats-io/stan.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

type (
	EventSvcCfg struct {
		l          log.Logger
		Publisher  *event.Publisher
		Subscriber *event.Subscriber
	}
	eventSvc struct {
		pub *event.Publisher
		sub *event.Subscriber
		l   *log.Logger
	}
)

func New(conf *EventSvcCfg) *eventSvc {
	return &eventSvc{
		conf.Publisher,
		conf.Subscriber,
		&conf.l,
	}
}

func (e *eventSvc) PrintEvent(ctx context.Context, event *pb.Event) (ept *empty.Empty, err error) {
	if event.GetEventData() == "" {
		return ept, errors.New("received empty event data")
	}
	e.l.Printf("Received event data: %s", event.EventData)
	return ept, nil
}

func (e *eventSvc) ListenNats(ctx context.Context) {
	mcb := func(msg *stan.Msg) {
		e.l.Printf("Received event: %s", string(msg.Data))
		var evt pb.Event
		if err := proto.Unmarshal(msg.Data, &evt); err != nil {
			e.l.Printf("Error unmarshal NATS event to proto: %v", err)
		}
		_, err := e.PrintEvent(ctx, &evt)
		if err != nil {
			e.l.Printf("Error while calling PrintEvent method: %v", err)
		}
	}
	if err := e.sub.Subscribe("test", "all-test", mcb); err != nil {
		e.l.Printf("Error subscribing to channel : %v", err)
	}
}

func (e *eventSvc) Run() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		e.l.Print("Subscribing to NATS/STAN")
		e.ListenNats(ctx)
	}()

	grpcAddr := "http://127.0.0.1:3333"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		e.l.Fatalf("Failed to listen to address / port: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterEventServiceServer(grpcServer, e)

	e.l.Printf("Serving grpc server on: %s", grpcAddr)
	err = grpcServer.Serve(lis)
	if err != nil {
		e.l.Fatalf("Failure serving grpc service: %v", err)
	}
}
