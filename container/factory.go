package container

import (
	"reflect"
)

type ObjectFactory interface {
	// Return an instance (possibly shared or independent) of the object managed by this factory.
	GetInstance() (interface{}, error)
}

type ObjectFactoryMethod func() (interface{}, error)

func (m ObjectFactoryMethod) GetInstance() (interface{}, error) { return m() }

// The root interface for accessing a container.
type InstanceFactory interface {
	// Does this factory contain a definition or externally registered singleton instance with the given name?
	ContainsInstance(name string) bool

	// Return the aliases for the given name, if any.
	GetAliases(name string) []string

	// Return an instance, which may be shared or independent, of the specified instance.
	GetInstanceByType(t reflect.Type, args ...interface{}) interface{}

	// Return an instance, which may be shared or independent, of the specified instance.
	GetInstanceByName(name string, args ...interface{}) interface{}

	// Determine the type of the bean with the given name.
	GetType(name string) reflect.Type

	// Is this a prototype? That is, will getInstanceByName(string) always return independent instances?
	IsPrototype(name string) bool

	// Is this a shared singleton? That is, will getInstanceByName(string) always return the same instance?
	IsSingleton(name string) bool
}

type HierarchicalInstanceFactory interface {
	InstanceFactory

	// Return whether the local instance factory contains a instance of the given name, ignoring instances defined in ancestor contexts.
	ContainsLocalInstance(name string) bool

	// Return the parent instance factory, or nil if there is none.
	Parent() InstanceFactory
}
