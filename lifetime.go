package di

type Lifetime int

const (
	LifetimeSingleton Lifetime = iota
	LifetimePerDependency
	LifetimePerHttpRequest
)

var lifetimes = map[Lifetime]bool{
	LifetimeSingleton:      true,
	LifetimePerDependency:  true,
	LifetimePerHttpRequest: true,
}
