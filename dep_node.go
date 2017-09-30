package di

import "reflect"

type depNode struct {
	Constructor reflect.Value
	DependsOn   []reflect.Type
	Edges       map[reflect.Type]*depNode
	Lifetime    Lifetime
	ReturnsErr  bool
	Type        reflect.Type
}

func newDepNode(constructor reflect.Value, lifetime Lifetime) *depNode {
	var node depNode

	node.Constructor = constructor
	node.Lifetime = lifetime
	node.Type = constructor.Out(0)

	if constructor.NumOut() == 2 {
		node.ReturnsErr = true
	}

	numIn := constructor.NumIn()
	deps := make([]reflect.Type, numIn)

	for index := range deps {
		deps[index] = constructor.In(index)
	}

	node.DependsOn = deps
	node.Edges = make(map[reflect.Type]*depNode, numIn)

	return &node
}

func (dn *depNode) AddEdge(node *depNode) {
	for _, dependsOn := range dn.DependsOn {
		if dependsOn == node.Type {
			dn.Edges[dependsOn] = node
			return
		}
	}
}

func (dn *depNode) MissingDependencies() []reflect.Type {
	missing := make([]reflect.Type, 0, len(dn.DependsOn))

	for _, dependsOn := range dn.DependsOn {
		_, hasEdge := dn.Edges[dependsOn]

		if hasEdge == false {
			missing = append(missing, dependsOn)
		}
	}

	return missing
}
