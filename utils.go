package dicontainer

import (
	"fmt"
	"reflect"
)

// setWithoutArgs sets dependencies without arguments.
func (d *di) separateCreators(creators ...any) ([]any, []any, error) {
	withArgs := make([]any, 0, len(creators))
	withoutArgs := make([]any, 0, len(creators))

	for _, creator := range creators {
		creatorType := reflect.TypeOf(creator)
		if creatorType.Kind() != reflect.Func {
			return nil, nil, fmt.Errorf("invalid creator type: %v", creator)
		}

		if creatorType.NumOut() == 0 {
			return nil, nil, fmt.Errorf("invalid creator, no return value: %v", creator)
		}

		if creatorType.NumIn() == 0 {
			withoutArgs = append(withoutArgs, creator)
		} else {
			withArgs = append(withArgs, creator)
		}
	}

	return withArgs, withoutArgs, nil
}

// setWithoutArgs sets dependencies without arguments.
func (d *di) setWithoutArgs(creators ...any) error {
	for _, creator := range creators {
		creatorType := reflect.TypeOf(creator)
		creatorFunc := reflect.ValueOf(creator)
		dependency := creatorType.Out(0)

		if dependency.Kind() != reflect.Ptr &&
			dependency.Kind() != reflect.Struct &&
			dependency.Kind() != reflect.Interface {
			return fmt.Errorf("invalid dependency type: %v", dependency)
		}

		depName := dependency.Name()
		if dependency.Kind() == reflect.Ptr {
			depName = dependency.Elem().Name()
		}

		if depName == "" {
			return fmt.Errorf("undefined dependency interface: %v", dependency)
		}

		// Call creator without arguments.
		result := creatorFunc.Call(nil)

		if err := d.setResult(depName, result); err != nil {
			return err
		}
	}

	return nil
}

// setWithArgs sets dependencies with arguments.
func (d *di) setWithArgs(creators ...any) error {
	for _, creator := range creators {
		creatorType := reflect.TypeOf(creator)
		dependency := creatorType.Out(0)

		if dependency.Kind() != reflect.Ptr &&
			dependency.Kind() != reflect.Struct &&
			dependency.Kind() != reflect.Interface {
			return fmt.Errorf("invalid dependency type: %v", dependency)
		}

		depName := dependency.Name()
		if dependency.Kind() == reflect.Ptr {
			depName = dependency.Elem().Name()
		}

		if depName == "" {
			return fmt.Errorf("undefined dependency interface: %v", dependency)
		}

		// Get arguments from creator.
		numIn := creatorType.NumIn()
		args := make([]reflect.Value, numIn)
		for i := 0; i < numIn; i++ {
			argType := creatorType.In(i)

			argName := argType.Name()
			if argType.Kind() == reflect.Ptr {
				argName = argType.Elem().Name()
			}

			arg, ok := d.services[argName]
			if !ok {
				return fmt.Errorf("missing dependency: %s", argName)
			}

			args[i] = arg
		}

		// Call creator with arguments.
		results := reflect.ValueOf(creator).Call(args)

		if err := d.setResult(depName, results); err != nil {
			return err
		}
	}

	return nil
}

// setResult sets result.
// The result of the creator execution must be a structure object or a pointer to a structure object.
// Optionally, the second argument returned must be an error.
func (d *di) setResult(key string, results []reflect.Value) error {
	if len(results) == 0 || len(results) > 2 {
		return fmt.Errorf("invalid result, expected 1 or 2 arguments: %v", results)
	}

	if len(results) > 1 {
		errType := reflect.TypeOf((*error)(nil)).Elem()
		if !results[1].Type().AssignableTo(errType) {
			return fmt.Errorf("invalid result, expected 2nd argument to be an error: %v", results)
		}

		if !results[1].IsNil() {
			return results[1].Interface().(error)
		}
	}

	d.services[key] = results[0]

	return nil
}
