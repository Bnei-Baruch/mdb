package events

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
	"github.com/spf13/viper"
)

var eventHandlers []EventHandler

func InitEmitter() (*BufferedEmitter, error) {
	// Setup events handlers
	eventHandlers = make([]EventHandler, 0)
	hNames := viper.GetStringSlice("events.handlers")
	if len(hNames) > 0 {
		for i := range hNames {
			switch hNames[i] {
			case "logger":
				eventHandlers = append(eventHandlers, new(LoggerEventHandler))
			case "nats":
				log.Info("Initializing nats streaming event handler")
				h, err := NewNatsStreamingEventHandler(
					viper.GetString("nats.subject"),
					viper.GetString("nats.cluster-id"),
					viper.GetString("nats.client-id"),
					stan.NatsURL(viper.GetString("nats.url")),
					stan.PubAckWait(viper.GetDuration("nats.pub-ack-wait")),
				)
				if err != nil {
					log.Errorf("Error connecting to nats streaming server: %s", err)
				} else {
					eventHandlers = append(eventHandlers, h)
				}
			default:
				log.Warnf("Unknown event handler: %s", hNames[i])
			}
		}
	}

	return NewBufferedEmitter(viper.GetInt("events.emitter-size"), eventHandlers...)
}

func CloseEmitter() {
	log.Infof("Closing event handlers")
	for i := range eventHandlers {
		if h, ok := eventHandlers[i].(io.Closer); ok {
			if err := h.Close(); err != nil {
				log.Fatal("Close event handler:", err)
			}
		}
	}
}
