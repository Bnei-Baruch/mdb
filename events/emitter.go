package events

import (
	"io"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

type EventEmitter interface {
	Emit(...Event)
}

type NoopEmitter struct{}

func (e *NoopEmitter) Emit(event ...Event) {}

type BufferedEmitter struct {
	running  bool
	queue    chan Event
	entropy  io.Reader
	handlers []EventHandler
}

func NewBufferedEmitter(size int, handlers ...EventHandler) (*BufferedEmitter, error) {
	if len(handlers) == 0 {
		return nil, errors.New("At least one handler is required")
	}

	emitter := new(BufferedEmitter)
	emitter.init(size)
	emitter.handlers = handlers

	return emitter, nil
}

func (e *BufferedEmitter) Emit(events ...Event) {
	if e.running && len(events) > 0 {
		for i := range events {
			e.queue <- events[i]
		}
	}
}

func (e *BufferedEmitter) Shutdown() {
	close(e.queue)
}

func (e *BufferedEmitter) init(size int) {
	e.queue = make(chan Event, size)
	e.entropy = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	go e.run()
}

func (e *BufferedEmitter) run() {
	e.running = true

	for event := range e.queue {
		event.ID = ulid.MustNew(ulid.Now(), e.entropy).String()
		for i := range e.handlers {
			e.handlers[i].Handle(event)
		}
	}

	e.running = false
}
