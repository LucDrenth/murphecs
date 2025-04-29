package app

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/log"
)

type System any

type systemEntry struct {
	system reflect.Value
	params []reflect.Value
}

func (s *systemEntry) exec(logger log.Logger) {
	result := s.system.Call(s.params)

	if len(result) == 1 {
		returnedError, isErr := result[0].Interface().(error)
		if isErr {
			logger.Error(fmt.Sprintf("system returned error: %v\n", returnedError))
		}
	}
}

type SystemSet struct {
	systems []systemEntry
}

func (s *SystemSet) add(sys System, world *ecs.World, logger log.Logger, resources *resourceStorage) error {
	systemValue := reflect.ValueOf(sys)
	queryType := reflect.TypeOf((*ecs.Query)(nil)).Elem()

	if err := validateSystem(systemValue, queryType, resources); err != nil {
		return fmt.Errorf("failed to validate system: %w", err)
	}

	numberOfParams := systemValue.Type().NumIn()
	params := make([]reflect.Value, numberOfParams)

	for i := range numberOfParams {
		parameterType := systemValue.Type().In(i)

		if parameterType.Implements(queryType) {
			query, ok := reflect.New(parameterType.Elem()).Interface().(ecs.Query)
			if !ok {
				return fmt.Errorf("failed to cast param to query")
			}

			err := query.PrepareOptions()
			if err != nil {
				return fmt.Errorf("failed to prepare query param: %w", err)
			}

			params[i] = reflect.ValueOf(query)
		} else if parameterType == reflect.TypeFor[*ecs.World]() {
			params[i] = reflect.ValueOf(world)
		} else if parameterType == reflect.TypeFor[log.Logger]() {
			params[i] = reflect.ValueOf(logger)
		} else {
			resource, err := resources.getReflectResource(parameterType)
			if err != nil {
				return fmt.Errorf("received unexpected system parameter: %w", err)
			}

			params[i] = resource
		}
	}

	entry := systemEntry{system: systemValue, params: params}
	s.systems = append(s.systems, entry)
	return nil
}

func validateSystem(sys reflect.Value, queryType reflect.Type, resources *resourceStorage) error {
	if sys.Kind() != reflect.Func {
		return errors.New("not a function")
	}

	if err := validateSystemReturnTypes(sys); err != nil {
		return fmt.Errorf("invalid return type(s): %w", err)
	}

	if err := validateSystemParameters(sys, queryType, resources); err != nil {
		return fmt.Errorf("invalid parameter(s): %w", err)
	}

	return nil
}

func validateSystemReturnTypes(systemValue reflect.Value) error {
	numberOfSystemReturnValues := systemValue.Type().NumOut()

	if numberOfSystemReturnValues == 0 {
		return nil
	}

	if numberOfSystemReturnValues == 1 {
		returnType := systemValue.Type().Out(0)

		if returnType == reflect.TypeFor[error]() {
			return nil
		}

		return fmt.Errorf("return type is %s but must be either error or nothing", returnType.String())
	}

	return fmt.Errorf("has %d return values but must have either 1 (error) or 0", numberOfSystemReturnValues)
}

func validateSystemParameters(systemValue reflect.Value, queryType reflect.Type, resources *resourceStorage) error {
	for i := range systemValue.Type().NumIn() {
		parameterType := systemValue.Type().In(i)

		if parameterType.Implements(queryType) {
			if err := validateQueryParameter(parameterType); err != nil {
				return fmt.Errorf("query parameter %d of type %s is not valid: %w", i, parameterType.String(), err)
			}
		} else if parameterType == reflect.TypeFor[*ecs.World]() {
			return nil
		} else if parameterType == reflect.TypeFor[log.Logger]() {
			return nil
		} else {
			_, err := resources.getReflectResource(parameterType)
			if err != nil {
				return fmt.Errorf("system parameter %d of type %s is not valid: %w", i+1, parameterType.String(), err)
			}

			return nil
		}
	}

	return nil
}

func validateQueryParameter(_ reflect.Type) error {
	// Golang does not have reflect for generics (as of Go 1.24). Thus we can not really
	// implement this yet.
	//
	// Once Golang gets generics reflection, here are some of the things we'd want to check:
	// - Return error if duplicates in the components to query
	// - Return error if marking a component as optional while that component is not queried
	// - Return error if there are any duplicate filters

	return nil
}
