package event

import (
	"github.com/hashicorp/go-uuid"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
)

type Subscriber struct {
	clusterID string
	id        string
	l         *log.Logger
	m         sync.Mutex
	conn      stan.Conn
	sub       stan.Subscription
}

func NewSubscriber(c *Conf) *Subscriber {
	return &Subscriber{
		clusterID: c.ClusterID,
		id:        c.ClientID,
	}
}

func (s *Subscriber) Connect(opts ...stan.Option) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	newID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	opts = append(opts, stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		s.l.Printf("Connection to NATS server lost: %v", err)
	}))
	if s.conn, err = stan.Connect(s.clusterID,
		s.id+"-"+newID, opts...); err != nil {
		s.l.Printf("nats connectivity status is DISCONNECTED: %v", err)
		return err
	}
	s.l.Printf("NATS Listening to: CID: %s", s.clusterID)
	return nil
}

// Subscribe takes subject to subscribe to
// and message callback function
func (s *Subscriber) Subscribe(subject, durableName string, mcb func(*stan.Msg)) error {
	var err error
	qgroup := "reports"
	s.sub, err = s.conn.QueueSubscribe(subject, qgroup, mcb, stan.StartWithLastReceived(),
		stan.DurableName(durableName))
	if err != nil {
		return err
	}
	return nil
}

func (s *Subscriber) Unsubscribe() {
	s.sub.Unsubscribe()
}
