package event

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	uuid "github.com/hashicorp/go-uuid"
	"sync"
	"time"

	"github.com/nats-io/stan.go"
)

type (
	Conf struct {
		ClusterID string
		Log *log.Logger
		ID        string
		ClientID  string
	}

	Publisher struct {
		clusterID string
		id        string
		l *log.Logger
		m         sync.Mutex
		conn      stan.Conn
		acb       func(string, error)
	}
	
	Event struct {
		ServiceId string `json:"serviceId"`
		Channel string `json:"channel"`
		AggregateId string `json:"aggregateId"`
		AggregateType string `json:"aggregateType"`
		EventId string `json:"eventId"`
		EventType string `json:"eventType"`
		EventData []byte `json:"eventData"`
		Originator string `json:"originator"`
		CreatedAt string `json:"createdAt"`
	}
)

func New(c *Conf) *Publisher {
	return &Publisher{
		clusterID: c.ClusterID,
		l: c.Log,
		id:        c.ID,
	}
}

func (p *Publisher) Connect(opts ...stan.Option) (err error) {

	p.m.Lock()
	defer p.m.Unlock()
	newID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	if p.conn, err = stan.Connect(p.clusterID , p.id + "-" + newID, opts...); err != nil {
		return nil
	}
	return nil
}

func (p *Publisher) Publish(evt *Event) error {
	eventID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	e := &Event{
		ServiceId:     evt.ServiceId,
		Channel:       evt.Channel,
		AggregateId:   evt.AggregateId,
		AggregateType: evt.AggregateType,
		EventId:       eventID,
		EventType:     evt.EventType,
		EventData:     evt.EventData,
		Originator:    evt.Originator,
		CreatedAt:     evt.CreatedAt,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	var guid string
	ch := make(chan bool, 1)
	p.acb = func(lguid string, err error) {
		p.m.Lock()
		p.l.Printf("Received ACK: %s", lguid)
		defer p.m.Unlock()
		if err != nil {
			p.l.Printf("Error in server ack: %s => %v",
				guid, err)
		}
		if lguid != guid {
			p.l.Printf("Error expecting matching guid in ack: %s wanted %s",
				lguid, guid)
		}
		ch <- true
	}
	p.m.Lock()
	guid, err = p.conn.PublishAsync(evt.Channel, b, p.acb)
	if err != nil {
		p.l.Printf("Error during async publish: %v", err)
	}
	p.m.Unlock()
	if guid == "" {
		p.l.Print("Expected non-empty guid to be returned.")
	}
	p.l.Printf(
		"Published: Event: %s\t Service: %s, Origin: %s, Data: %s\n",
		e.EventId,
		e.ServiceId,
		e.Originator,
		e.EventData,
	)
	select {
	case <-ch:
		return nil
	case <-time.After(10 * time.Second):
		return status.Error(codes.Canceled, "timeout")
	}
}