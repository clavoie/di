package di

type Lifetime int

const (
	LifetimeSingleton Lifetime = iota
	LifetimePerDependency
	LifetimePerHttpRequest
	LifetimePerResolution
)

var lifetimes = map[Lifetime]bool{
	LifetimeSingleton:      true,
	LifetimePerDependency:  true,
	LifetimePerHttpRequest: true,
	LifetimePerResolution:  true,
}
