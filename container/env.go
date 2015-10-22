package container

type Environment interface {
}

type ConfigurableEnvironment interface {
	Environment
}

type EnvironmentCapable interface {
	Environment() Environment
}

type AbstractEnvironment struct {
	ConfigurableEnvironment
}
