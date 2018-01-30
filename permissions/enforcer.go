package permissions

import (
	"github.com/casbin/casbin"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/mdb/bindata"
)

func NewEnforcer() (*casbin.Enforcer, error) {
	e := casbin.NewEnforcer()
	e.EnableLog(false)

	// load model
	pModel, err := bindata.Asset("data/permissions_model.conf")
	if err != nil {
		return nil, errors.Wrap(err, "Load permissions_model.conf")
	}
	e.SetModel(casbin.NewModel(string(pModel)))

	e.InitWithModelAndAdapter(casbin.NewModel(string(pModel)), NewBindataPolicyAdapter())
	return e, nil
}
