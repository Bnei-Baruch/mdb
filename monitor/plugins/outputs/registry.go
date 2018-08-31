package outputs

import "github.com/Bnei-Baruch/mdb/monitor/interfaces"

type Creator func() interfaces.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
