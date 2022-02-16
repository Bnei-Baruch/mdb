package permissions

import (
	_ "embed"

	"github.com/casbin/casbin"
)

//go:embed permissions_model.conf
var permissionsModel string

func NewEnforcer() *casbin.Enforcer {
	e := casbin.NewEnforcer()
	e.EnableLog(false)
	e.SetModel(casbin.NewModel(permissionsModel))
	e.InitWithModelAndAdapter(casbin.NewModel(permissionsModel), NewEmbeddedPolicyAdapter())
	return e
}
