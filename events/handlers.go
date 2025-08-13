package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type EventHandler interface {
	Handle(Event)
	Close(context.Context) error
}

type LoggerEventHandler struct{}

func (eh *LoggerEventHandler) Handle(event Event) {
	log.WithFields(log.Fields{
		"id":      event.ID,
		"type":    event.Type,
		"rloc":    event.ReplicationLocation,
		"payload": event.Payload,
	}).Info("event")
}

func (eh *LoggerEventHandler) Close(ctx context.Context) error {
	return nil
}

// Nats
type NatsEventHandler struct {
	nc       *nats.Conn
	js       jetstream.JetStream
	ncClosed chan struct{}
}

func (eh *NatsEventHandler) Handle(event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Errorf("NatsEvetnHandler.Handle: json.Marshal: %s %v", event.ID, err)
		return
	}

	subject := fmt.Sprintf("mdb.%s", strings.ToLower(event.Type))

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err = eh.js.Publish(ctx, subject, payload, jetstream.WithRetryAttempts(3))
	if err != nil {
		log.Errorf("NatsEvetnHandler.Handle: jetstream.Publish: %s %v", event.ID, err)
		return
	}
}

func (eh *NatsEventHandler) Close(ctx context.Context) error {
	if err := eh.nc.Drain(); err != nil {
		return fmt.Errorf("NatsEventhandler: nc.Drain(): %w", err)
	}

	select {
	case <-eh.ncClosed:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("NatsEventHandler drain ctx.Done: %w", ctx.Err())
	}
}

func (eh *NatsEventHandler) closedCallback(conn *nats.Conn) {
	log.Info("Nats connection closed.")
	eh.ncClosed <- struct{}{}
}

func NewNatsStreamingEventHandler(natsURL string) (*NatsEventHandler, error) {
	eh := new(NatsEventHandler)
	eh.ncClosed = make(chan struct{})

	var err error
	log.Info("Initialize connection to nats")
	eh.nc, err = nats.Connect(natsURL, nats.ClosedHandler(eh.closedCallback))
	if err != nil {
		return nil, fmt.Errorf("nats.Connect: %w", err)
	}

	eh.js, err = jetstream.New(eh.nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream.New: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = eh.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "MDB",
		Description: "Events stream for MDB",
		Subjects:    []string{"mdb.*"},
		MaxMsgs:     4096,
		Storage:     jetstream.FileStorage,
	})
	if err != nil {
		return nil, fmt.Errorf("jetstream.CreateOrUpdateStream: %w", err)
	}

	return eh, nil
}
