package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type HandlersSuite struct {
	suite.Suite
}

func (suite *HandlersSuite) SetupSuite() {
	common.InitConfig()
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

func (suite *HandlersSuite) TestNatsHandler() {
	natsUrl := viper.GetString("nats.url")
	fmt.Println("Initializing test NATS: ", natsUrl)
	handler, err := NewNatsStreamingEventHandler(natsUrl)
	suite.Require().Nil(err)

	// Handle 100 events.
	for i := 0; i < 100; i++ {
		handler.Handle(Event{ID: fmt.Sprintf("test-event-%d", i)})
	}

	time.Sleep(20 * time.Millisecond)

	// Close handler (before it complete publishing all 100 events).
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	suite.Require().Nil(handler.Close(ctx))
	<-ctx.Done()
}
