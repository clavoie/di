package di

import "reflect"

type singleton struct {
	Node  *depNode
	Value reflect.Value
}

func newSingleton(node *depNode) *singleton {
	return &singleton{
		Node: node,
	}
}
