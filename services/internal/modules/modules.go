// Package modules is the single place where feature modules are registered.
package modules

import (
	"errors"

	"github.com/hungphan1911/tapestry/services/internal/core"
	"github.com/hungphan1911/tapestry/services/internal/modules/finance"
)

// Migrations lists every module's migrations, in the order they are applied.
func Migrations() []core.Migration {
	return []core.Migration{
		finance.Migration(),
	}
}

// Build constructs all modules. If one fails, those already opened are closed.
func Build(pg core.PostgresConfig) ([]core.Module, error) {
	constructors := []func(core.PostgresConfig) (core.Module, error){
		finance.New,
	}

	var built []core.Module
	for _, newModule := range constructors {
		m, err := newModule(pg)
		if err != nil {
			err = errors.Join(err, CloseAll(built))
			return nil, err
		}
		built = append(built, m)
	}
	return built, nil
}

func CloseAll(mods []core.Module) error {
	var errs []error
	for _, m := range mods {
		errs = append(errs, m.Close())
	}
	return errors.Join(errs...)
}
