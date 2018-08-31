package discard

import (
	"github.com/Bnei-Baruch/mdb/monitor/interfaces"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/outputs"
)

type Discard struct{}

func (d *Discard) TryParseConfigurations(outputConfigs map[string]interface{}) error { return nil }
func (d *Discard) Connect() error                                                    { return nil }
func (d *Discard) Close() error                                                      { return nil }
func (d *Discard) SampleConfig() string                                              { return "" }
func (d *Discard) Description() string                                               { return "Send metrics to nowhere at all" }
func (d *Discard) Write(metrics []interfaces.Metric) error                           { return nil }

func init() {
	outputs.Add("discard", func() interfaces.Output { return &Discard{} })
}
