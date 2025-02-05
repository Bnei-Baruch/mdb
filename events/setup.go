package events

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"
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
				log.Info("Initializing nats jetstreaming event handler")
				if h, err := NewNatsStreamingEventHandler(viper.GetString("nats.url")); err != nil {
					// Fail if NATS not starting.
					return nil, fmt.Errorf("Error connecting to nats streaming server: %w", err)
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

func CloseEmitter(ctx context.Context) {
	log.Infof("Closing event handlers")
	for i := range eventHandlers {
		if err := eventHandlers[i].Close(ctx); err != nil {
			log.Error("Close event handler:", err)
		}
	}
}
