package events

import (
	"context"
	"encoding/json"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
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

const queueSize = 4096
const nats_tmp_file = "__DO_NOT_REMOVE_nats-event-handler.tmp"

type NatsStreamingEventHandler struct {
	sc      stan.Conn
	subject string
	ch      chan *Event
	stopCH  chan bool
}

func NewNatsStreamingEventHandler(subject, clusterID, clientID string,
	options ...stan.Option) (*NatsStreamingEventHandler, error) {
	eh := new(NatsStreamingEventHandler)

	// connect to nats

	// Unfortunately, there is an open issue regarding connection failures on startup.
	// see https://github.com/nats-io/go-nats/issues/195
	// we should upgrade as soon as it's fixed !
	var err error
	eh.sc, err = stan.Connect(clusterID, clientID, options...)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}

	eh.subject = subject
	eh.stopCH = make(chan bool)
	eh.ch = make(chan *Event, queueSize)

	// start runner
	go eh.run()

	// try to load unprocessed events from temp file (last shutdown)
	if err := eh.loadFromFile(); err != nil {
		return nil, errors.Wrap(err, "load from file")
	}

	return eh, nil
}

func (eh *NatsStreamingEventHandler) Handle(event Event) {
	if len(eh.ch) < cap(eh.ch) {
		eh.ch <- &event
	} else {
		log.Warnf("nats: buffer limit reached, dropping event %s", event.ID)
	}
}

func (eh *NatsStreamingEventHandler) Close(ctx context.Context) error {
	// whatever happens, close connection to nats (second time is noop)
	defer eh.sc.Close()

	// close channel
	log.Infof("nats: close events channel")
	close(eh.ch)

	// poll timer until context timeout or no more messages in queue
	log.Infof("nats: drain %d events in queue", len(eh.ch))
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		if len(eh.ch) == 0 {
			break
		}
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
		}
	}

	// stop runner
	log.Infof("nats: stop runner")
	eh.stopCH <- true

	// drain what we have left in queue to file
	if err := eh.drainToFile(); err != nil {
		return errors.Wrap(err, "drain to file")
	}

	// close connection to nats
	log.Infof("nats: close connection")
	return eh.sc.Close()
}

func (eh *NatsStreamingEventHandler) run() {
	for {
		select {
		case <-eh.stopCH:
			return
		case event := <-eh.ch:
			if event != nil {
				if err := eh.publish(event); err != nil {
					log.Errorf("nats: publish error %s", err.Error())
				}
			}
		}
	}
}

func (eh *NatsStreamingEventHandler) publish(event *Event) error {
	log.Infof("nats: publish event %s", event.ID)

	b, err := json.Marshal(event)
	if err != nil {
		log.Errorf("nats: json.Marshal event [%s]: %s", event.ID, err.Error())
		return nil // not a nats related error. report don't choke
	}

	// sync publish, timeout is set on the nats client
	err = eh.sc.Publish(eh.subject, b)
	if err != nil {
		return errors.Wrapf(err, "publish event [%s]", event.ID)
	}

	return nil
}

func (eh *NatsStreamingEventHandler) loadFromFile() error {
	f, err := os.Open(nats_tmp_file)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "os.Open")
	}

	defer func() {
		f.Close()
		os.Remove(nats_tmp_file)
	}()

	if f != nil {
		var evnts []*Event
		if err := json.NewDecoder(f).Decode(&evnts); err != nil {
			return errors.Wrap(err, "json.Decode")
		}

		if len(evnts) > 0 {
			log.Warnf("nats: tmp file has %d events. queuing...", len(evnts))
			for i := range evnts {
				eh.ch <- evnts[i]
			}
		}
	}

	return nil
}

func (eh *NatsStreamingEventHandler) drainToFile() error {
	evnts := make([]*Event, 0)

	// drain channel
	// channel is expected to be closed by now
	for e := range eh.ch {
		evnts = append(evnts, e)
	}

	if len(evnts) > 0 {
		log.Infof("nats: drain %d unprocessed events to tmp file", len(evnts))

		f, err := os.Create(nats_tmp_file)
		if err != nil {
			return errors.Wrap(err, "create tmp file")
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(evnts); err != nil {
			return errors.Wrap(err, "json.Encode")
		}
	}

	return nil
}
