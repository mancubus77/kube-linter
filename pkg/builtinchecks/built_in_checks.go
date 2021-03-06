package builtinchecks

import (
	"sync"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
)

var (
	box = packr.NewBox("./yamls")

	loadOnce sync.Once
	list     []check.Check
	loadErr  error
)

// LoadInto loads built-in checks into the registry.
func LoadInto(registry checkregistry.CheckRegistry) error {
	checks, err := List()
	if err != nil {
		return err
	}
	for _, chk := range checks {
		if err := registry.Register(&chk); err != nil {
			return errors.Wrapf(err, "registering default check %s", chk.Name)
		}
	}
	return nil
}

// List lists built-in checks.
func List() ([]check.Check, error) {
	loadOnce.Do(func() {
		for _, fileName := range box.List() {
			contents, err := box.Find(fileName)
			if err != nil {
				loadErr = errors.Wrapf(err, "loading default check from %s", fileName)
				return
			}
			var chk check.Check
			if err := yaml.Unmarshal(contents, &chk); err != nil {
				loadErr = errors.Wrapf(err, "unmarshalling default check from %s", fileName)
				return
			}
			list = append(list, chk)
		}
	})
	return list, loadErr
}
