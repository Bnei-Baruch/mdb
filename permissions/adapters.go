package permissions

import (
	"bufio"
	"bytes"
	_ "embed"
	"io"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/pkg/errors"
)

//go:embed permissions_policy.csv
var permissionsPolicy []byte

type EmbeddedPolicyAdapter struct {
}

func NewEmbeddedPolicyAdapter() *EmbeddedPolicyAdapter {
	return new(EmbeddedPolicyAdapter)
}

// LoadPolicy loads all policy rules from the storage.
func (a *EmbeddedPolicyAdapter) LoadPolicy(model model.Model) error {
	buf := bufio.NewReader(bytes.NewReader(permissionsPolicy))
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
func (a *EmbeddedPolicyAdapter) SavePolicy(model.Model) error {
	return errors.New("not implemented")
}

// AddPolicy adds a policy rule to the storage.
func (a *EmbeddedPolicyAdapter) AddPolicy(string, string, []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *EmbeddedPolicyAdapter) RemovePolicy(string, string, []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *EmbeddedPolicyAdapter) RemoveFilteredPolicy(string, string, int, ...string) error {
	return errors.New("not implemented")
}
