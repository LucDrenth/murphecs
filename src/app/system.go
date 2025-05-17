package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type System any

type systemEntry struct {
	system reflect.Value
	params []reflect.Value
}

func (s *systemEntry) exec() error {
	result := s.system.Call(s.params)

	if len(result) == 1 {
		returnedError, isErr := result[0].Interface().(error)
		if isErr {
			return returnedError
		}
	}

	return nil
}

type SystemSet struct {
	systems            []systemEntry
	systemParamQueries []ecs.Query
}

func (s *SystemSet) Exec(world *ecs.World) []error {
	err := s.handleSystemParamQueries(world)
	if err != nil {
		return []error{
			fmt.Errorf("did not execute system set because query failed: %w", err),
		}
	}

	return s.execSystems()
}

func (s *SystemSet) handleSystemParamQueries(world *ecs.World) error {
	for i := range s.systemParamQueries {
		query := s.systemParamQueries[i]
		if query.IsLazy() {
			query.ClearResults()
		} else {
			err := query.Exec(world)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SystemSet) execSystems() []error {
	errors := []error{}

	for i := range s.systems {
		err := s.systems[i].exec()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s *SystemSet) add(sys System, world *ecs.World, logger Logger, resources *resourceStorage) error {
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
			query, err := parseQueryParam(parameterType, world, logger)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrSystemParamQueryNotValid, err)
			}

			s.systemParamQueries = append(s.systemParamQueries, query)
			params[i] = reflect.ValueOf(query)
		} else if parameterType == reflect.TypeFor[*ecs.World]() {
			params[i] = reflect.ValueOf(world)
		} else {
			resource, err := resources.getReflectResource(parameterType)
			if err != nil {
				return fmt.Errorf("received unexpected system parameter: %w", err)
			}

			if parameterType.Kind() == reflect.Pointer {
				params[i] = resource
			} else {
				params[i] = resource.Elem()
			}
		}
	}

	entry := systemEntry{
		system: systemValue,
		params: params,
	}
	s.systems = append(s.systems, entry)
	return nil
}

func parseQueryParam(parameterType reflect.Type, world *ecs.World, logger Logger) (ecs.Query, error) {
	if parameterType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("can not be an interface")
	}

	query, ok := reflect.New(parameterType.Elem()).Interface().(ecs.Query)
	if !ok {
		return nil, fmt.Errorf("failed to cast param to query")
	}

	err := query.Prepare(world)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query param: %w", err)
	}

	warning := query.Validate()
	if warning != nil {
		logger.Warn(fmt.Sprintf("query %s is not optimized: %v", parameterType.String(), warning))
	}

	return query, nil
}

func validateSystem(sys reflect.Value, queryType reflect.Type, resources *resourceStorage) error {
	if sys.Kind() != reflect.Func {
		return ErrSystemNotAFunction
	}

	if err := validateSystemReturnTypes(sys); err != nil {
		return fmt.Errorf("%w: %w", ErrSystemInvalidReturnType, err)
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
			// query will be validated at a later step because it first needs to get prepared with `Query.Prepare` before
			// it can be validated with `Query.Validate`.
			return nil
		} else if parameterType == reflect.TypeFor[*ecs.World]() {
			return nil
		} else if parameterType == reflect.TypeFor[ecs.World]() {
			// ecs.World may not be used by-value because:
			//	1. it is a potentially big object and copying it could give bad performance
			//	2. it is probably unintended and would cause unexpected behavior
			return fmt.Errorf("system parameter %d: %w", i+1, ErrSystemParamWorldNotAPointer)
		} else {
			_, err := resources.getReflectResource(parameterType)
			if err == nil {
				return nil
			}

			if parameterType.Kind() != reflect.Pointer && reflect.PointerTo(parameterType).Implements(reflect.TypeFor[ecs.Query]()) {
				return fmt.Errorf("system parameter %d: %w", i+1, ErrSystemParamQueryNotAPointer)
			}

			return fmt.Errorf("system parameter %d: %w", i+1, ErrSystemParamNotValid)
		}
	}

	return nil
}

// systemToDebugString returns a reflection string of the system but with shortened paths.
func systemToDebugString(system System) string {
	result := reflect.TypeOf(system).String()
	result = strings.ReplaceAll(result, "github.com/lucdrenth/murphecs/src/", "murphecs/")
	return result
}
