package permissions

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/mdb/bindata"
	"github.com/Bnei-Baruch/mdb/utils"
)

type BindataPolicyAdapter struct {
}

func NewBindataPolicyAdapter() *BindataPolicyAdapter {
	return new(BindataPolicyAdapter)
}

// LoadPolicy loads all policy rules from the storage.
func (a *BindataPolicyAdapter) LoadPolicy(model model.Model) error {
	pPoicy, err := bindata.Asset("data/permissions_policy.csv")
	utils.Must(err)

	buf := bufio.NewReader(bytes.NewReader(pPoicy))
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		persist.LoadPolicyLine(line, model)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// SavePolicy saves all policy rules to the storage.
func (a *BindataPolicyAdapter) SavePolicy(model model.Model) error {
	return errors.New("not implemented")
}

// AddPolicy adds a policy rule to the storage.
func (a *BindataPolicyAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *BindataPolicyAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *BindataPolicyAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
