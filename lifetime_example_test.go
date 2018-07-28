package di_test

import (
	"fmt"
	"net/http"

	"github.com/clavoie/di"
)

type Singleton interface {
	Value() int
}
type PerDependency interface {
	Value() int
}
type PerResolve interface {
	Value() int
}
type Impl struct{ value int }

func (li *Impl) Value() int { return li.value }

type Dependent interface{}
type DependentImpl struct {
	S1, S2 Singleton
	D1, D2 PerDependency
	R1, R2 PerResolve
}

func NewDependent(s1, s2 Singleton, d1, d2 PerDependency, r1, r2 PerResolve) Dependent {
	return &DependentImpl{s1, s2, d1, d2, r1, r2}
}

func ExampleLifetime() {
	counter := 0
	newImpl := func() *Impl {
		counter += 1
		return &Impl{counter}
	}

	newSingleton := func() Singleton { return (Singleton)(newImpl()) }
	newPerDependency := func() PerDependency { return (PerDependency)(newImpl()) }
	newPerResolve := func() PerResolve { return (PerResolve)(newImpl()) }

	deps := []*di.Def{
		di.NewDef(newSingleton, di.Singleton),
		di.NewDef(newPerDependency, di.PerDependency),
		di.NewDef(newPerResolve, di.PerResolve),
		di.NewDef(NewDependent, di.PerDependency),
	}

	errFn := func(er *di.ErrResolve, w http.ResponseWriter, r *http.Request) { panic(er) }
	resolver, err := di.NewResolver(errFn, deps)
	if err != nil {
		panic(err)
	}

	var dependent Dependent
	resolveErr := resolver.Resolve(&dependent)
	if resolveErr != nil {
		panic(resolveErr)
	}

	// first resolution
	impl := dependent.(*DependentImpl)
	fmt.Println(impl.S1 == impl.S2, impl.S1.Value(), impl.S2.Value())
	fmt.Println(impl.D1 == impl.D2, impl.D1.Value(), impl.D2.Value())
	fmt.Println(impl.R1 == impl.R2, impl.R1.Value(), impl.R2.Value())

	resolveErr = resolver.Resolve(&dependent)
	if resolveErr != nil {
		panic(resolveErr)
	}

	impl2 := dependent.(*DependentImpl)
	fmt.Println(impl2.S1 == impl2.S2, impl2.S1.Value(), impl2.S2.Value())
	fmt.Println(impl2.D1 == impl2.D2, impl2.D1.Value(), impl2.D2.Value())
	fmt.Println(impl2.R1 == impl2.R2, impl2.R1.Value(), impl2.R2.Value())

	fmt.Println(impl.S1 == impl2.S1)
	// Output: true 1 1
	// false 2 3
	// true 4 4
	// true 1 1
	// false 5 6
	// true 7 7
	// true
}
