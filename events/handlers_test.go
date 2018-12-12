package events

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats-streaming-server/server"
	"github.com/stretchr/testify/suite"
)

const (
	clusterName = "my-test-cluster"
	clientName  = "test-client"
)

type HandlersSuite struct {
	suite.Suite
	nss *server.StanServer
}

func (suite *HandlersSuite) SetupSuite() {
	var err error
	suite.nss, err = server.RunServer(clusterName)
	suite.Require().Nil(err)
}

func (suite *HandlersSuite) TearDownSuite() {
	suite.nss.Shutdown()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

func (suite *HandlersSuite) TestNatsHandler() {
	handler, err := NewNatsStreamingEventHandler("test-subject", clusterName, clientName)
	suite.Require().Nil(err)

	// handle 100 events
	for i := 0; i < 100; i++ {
		handler.Handle(Event{ID: fmt.Sprintf("test-event-%d", i)})
	}

	// close handler (before it complete publishing all 100 events)
	suite.Require().Nil(handler.Close())

	// drain temp file with unpublished events using another, dummy handler
	handler2 := &NatsStreamingEventHandler{
		ch: make(chan *Event, queueSize),
	}
	suite.Require().Nil(handler2.loadFromFile())
	suite.NotEmpty(handler2.ch)
}
