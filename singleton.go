package di

import "reflect"

type singleton struct {
	node  *depNode
	value reflect.Value
}

func newSingleton(node *depNode) *singleton {
	return &singleton{
		node: node,
	}
}

func newSingletonValue(value reflect.Value) *singleton {
	return &singleton{
		value: value,
	}
}

func (s *singleton) SetValue(ins []reflect.Value) (reflect.Value, error) {
	value, err := s.node.NewValue(ins)

	if err != nil {
		return value, err
	}

	s.value = value
	return value, nil
}

func (s *singleton) Value() (reflect.Value, bool) {
	return s.value, s.value.IsValid()
}
