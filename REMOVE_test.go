package di

import (
	"reflect"
	"testing"
)

type IA interface {
	A() int
}

type aImpl struct {
	a int
}

func newA() IA {
	return &aImpl{1}
}

func (a *aImpl) A() int { return a.a }

type IB interface {
	B() int
}

type bImpl struct {
	b int
}

func newB(a IA) IB {
	return &bImpl{a.A() + 1}
}

func (b *bImpl) B() int { return b.b }

func TestXXX(t *testing.T) {
	cw := NewContainerWriter()
	err := cw.Add(newA, PerDependency)

	if err != nil {
		t.Fatal(err)
	}

	cw.Add(newB, PerDependency)

	if err != nil {
		t.Fatal(err)
	}

	c, err := cw.Build()

	if err != nil {
		t.Fatal(err)
	}

	var b IB
	err = c.Resolve(&b)

	if err != nil {
		t.Fatal(err)
	}

	if b.B() != 2 {
		t.Fatal(b.B())
	}

	fn := func(i int, b IB) int {
		return i + b.B()
	}

	cfn, err := c.Curry(fn)

	if err != nil {
		t.Fatal(err)
	}

	castCfn, isType := cfn.(func(int) int)

	if isType == false {
		t.Fatal(reflect.TypeOf(cfn))
	}

	if castCfn(1) != 3 {
		t.Fatal(castCfn(1))
	}
}
