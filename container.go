package di

type IContainer interface {
	Invoke(fn interface{})
}
