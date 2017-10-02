package di

import "reflect"

var aCounter = 0
var aType = reflect.TypeOf((*A)(nil)).Elem()
var bType = reflect.TypeOf((*B)(nil)).Elem()
var eType = reflect.TypeOf((*E)(nil)).Elem()

type A interface {
	A() int
}

type aImpl struct {
	a int
}

func (ai *aImpl) A() int { return ai.a }

func NewA() A {
	aCounter += 1
	return &aImpl{aCounter}
}

type B interface {
	B() (int, int)
}

type bImpl struct {
	A1, A2 int
}

func (bi *bImpl) B() (int, int) {
	return bi.A1, bi.A2
}

func NewB(a1, a2 A) B {
	return &bImpl{a1.A(), a2.A()}
}

type C interface{}
type D interface{}
type E interface{}

func NewC(d D, e E) C {
	return new(struct{})
}

func NewD(e E) D {
	return new(struct{})
}

func NewE(c C) E {
	return new(struct{})
}
