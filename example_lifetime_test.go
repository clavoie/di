package di_test

import (
	"fmt"

	"github.com/clavoie/di"
)

type LifetimeDep interface {
	Value() int
}

var lifetimeCounter = 0

type LifetimeImpl struct {
	value int
}

func NewLifetimeDep() LifetimeDep {
	lifetimeCounter += 1
	return &LifetimeImpl{lifetimeCounter}
}

func (li *LifetimeImpl) Value() int { return li.value }

func Example() {
	defs := di.NewDefs()
	err := defs.Add(NewLifetimeDep, di.Singleton)

	if err != nil {
		panic(err)
	}

	resolver, err := di.NewResolver(defs)
	if err != nil {
		panic(err)
	}

	var life1, life2 LifetimeDep
	err = resolver.Resolve(&life1)
	if err != nil {
		panic(err)
	}

	err = resolver.Resolve(&life2)
	if err != nil {
		panic(err)
	}

	fmt.Println(life1.Value(), life2.Value())

	// Output: 1 1
}
