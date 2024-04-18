package dicontainer

import (
	"reflect"
)

// DI represents a dependency injection interface.
type DI interface {
	// Set sets dependencies.
	Set(creators ...any) error

	// Get gets a dependency.
	Get(key string) any
}

// di is the default implementation of DI.
type di struct {
	services map[string]reflect.Value
}

// New creates a new DI instance.
func New() DI {
	return &di{make(map[string]reflect.Value)}
}

// Set sets dependencies.
func (d *di) Set(creators ...any) error {
	// Separate dependencies with and without arguments.
	withArgs, withoutArgs, err := d.separateCreators(creators...)
	if err != nil {
		return err
	}

	// Set dependencies without arguments.
	if err = d.setWithoutArgs(withoutArgs...); err != nil {
		return err
	}

	// Set dependencies with arguments.
	if err = d.setWithArgs(withArgs...); err != nil {
		return err
	}

	return nil
}

// Get gets a dependency.
func (d *di) Get(key string) any {
	dep, ok := d.services[key]
	if !ok {
		return nil
	}

	return dep.Interface()
}
