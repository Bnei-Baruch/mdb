package simple

import "github.com/Bnei-Baruch/mdb/monitor/plugins/inputs"
import "github.com/Bnei-Baruch/mdb/monitor/interfaces"

type Simple struct {
	Ok bool
}

func (s *Simple) Description() string {
	return "a demo plugin"
}

func (d *Simple) TryParseConfigurations(inputConfigs map[string]interface{}) error { return nil }

func (s *Simple) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  ok = true
`
}

func (s *Simple) Gather(acc interfaces.Accumulator) error {
	if s.Ok {
		acc.AddFields("state", map[string]interface{}{"value": "pretty good"}, nil)
	} else {
		acc.AddFields("state", map[string]interface{}{"value": "not great"}, nil)
	}

	return nil
}

func init() {
	inputs.Add("simple", func() interfaces.Input { return &Simple{} })
}
