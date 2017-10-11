package di

import (
	"fmt"
	"reflect"
	"strings"
)

type depNode struct {
	Constructor reflect.Value
	DependsOn   []reflect.Type
	Edges       map[reflect.Type]*depNode
	Lifetime    Lifetime
	ReturnsErr  bool
	Type        reflect.Type
	TypeName    string
}

func newDepNode(constructor reflect.Value, lifetime Lifetime, depMap map[reflect.Type]*depNode) *depNode {
	var node depNode

	node.Constructor = constructor
	node.Lifetime = lifetime

	constructorType := constructor.Type()
	node.Type = constructorType.Out(0)
	node.TypeName = node.Type.String()

	if constructorType.NumOut() == 2 {
		node.ReturnsErr = true
	}

	numIn := constructorType.NumIn()
	deps := make([]reflect.Type, numIn)
	edges := make(map[reflect.Type]*depNode, numIn)

	for index := range deps {
		inType := constructorType.In(index)
		deps[index] = inType

		edgeNode, hasNode := depMap[inType]
		if hasNode {
			edges[inType] = edgeNode
		}
	}

	node.DependsOn = deps
	node.Edges = edges

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

func (dn *depNode) CheckForCycle(seen []*depNode, checked map[*depNode]bool) error {
	hasSeen := func(dn2 *depNode) bool {
		for _, node := range seen {
			if node == dn2 {
				return true
			}
		}

		return false
	}

	for _, node := range dn.Edges {
		if checked[node] {
			continue
		}

		if hasSeen(node) {
			// print and return err
			path := make([]string, len(seen), len(seen)+1)
			for index, seenNode := range seen {
				path[index] = seenNode.Type.String()
			}
			path = append(path, node.Type.String())
			pathStr := strings.Join(path, "->")
			return fmt.Errorf("di: circular dependency detected: %v", pathStr)
		}

		seenCopy := make([]*depNode, len(seen), len(seen)+1)
		copy(seenCopy, seen)
		seenCopy = append(seenCopy, node)
		err := node.CheckForCycle(seenCopy, checked)

		if err != nil {
			return err
		}
	}

	checked[dn] = true
	return nil
}

func (dn *depNode) IsLeaf() bool {
	return len(dn.DependsOn) == 0
}

func (dn *depNode) NewValue(ins []reflect.Value) (reflect.Value, error) {
	outs := dn.Constructor.Call(ins)

	var err error
	if dn.ReturnsErr {
		val := outs[1].Interface()

		if val != nil {
			err = val.(error)
		}
	}

	return outs[0], err
}
