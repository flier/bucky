package container

import (
	"reflect"
	"strings"
)

const (
	FactoryBeanPrefix = "&"
)

type ObjectFactory interface {
	// Return an instance (possibly shared or independent) of the object managed by this factory.
	GetInstance() (interface{}, error)
}

type ObjectFactoryMethod func() (interface{}, error)

func (m ObjectFactoryMethod) GetInstance() (interface{}, error) { return m() }

type FactoryInstance interface {
	GetInstance() (interface{}, error)

	// Return the type of object that this FactoryInstance creates, or nil if not known in advance.
	ObjectType() reflect.Type

	// Is the object managed by this factory a singleton?
	IsSingleton() bool
}

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

	// Determine the type of the instance with the given name.
	GetType(name string) reflect.Type

	// Is this a prototype? That is, will getInstanceByName(string) always return independent instances?
	IsPrototype(name string) bool

	// Is this a shared singleton? That is, will getInstanceByName(string) always return the same instance?
	IsSingleton(name string) bool
}

type AutowireCapableInstanceFactory interface {
}

type ListableInstanceFactory interface {
	// Check if this instance factory contains a instance definition with the given name.
	ContainsInstanceDefinition(name string) bool

	// Return the registered InstanceDefinition for the specified instance
	GetInstanceDefinition(name string) InstanceDefinition

	// Return the number of instances defined in the factory.
	InstanceDefinitionCount() int

	// Return the names of all instances defined in this factory.
	InstanceDefinitionNames() []string

	// Return the names of instances matching the given type (including subclasses),
	// judging from either instance definitions or the value of ObjectType() in the case of FactoryInstances.
	GetInstanceNamesForType(t reflect.Type) []string

	// Return the instances that match the given object type (including subclasses),
	// judging from either instance definitions or the value of ObjectType() in the case of FactoryInstances.
	GetInstancesOfType(t reflect.Type) map[string]interface{}
}

type HierarchicalInstanceFactory interface {
	InstanceFactory

	// Return whether the local instance factory contains a instance of the given name, ignoring instances defined in ancestor contexts.
	ContainsLocalInstance(name string) bool

	// Return the parent instance factory, or nil if there is none.
	ParentInstanceFactory() InstanceFactory
}

type ConfigurableInstanceFactory interface {
	HierarchicalInstanceFactory
	SingletonInstanceRegistry

	// Determine whether the bean with the given name is a FactoryInstance.
	IsFactoryInstance(name string) (bool, error)
}

type ConfigurableListableInstanceFactory interface {
	ListableInstanceFactory
	AutowireCapableInstanceFactory
	ConfigurableInstanceFactory
}

type InstanceFactoryUtils struct {
}

func (u *InstanceFactoryUtils) transformedBeanName(name string) string {
	for {
		if !strings.HasPrefix(name, FactoryBeanPrefix) {
			break
		}

		name = name[len(FactoryBeanPrefix):]
	}

	return name
}

func (u *InstanceFactoryUtils) isFactoryDereference(name string) bool {
	return strings.HasPrefix(name, FactoryBeanPrefix)
}

type FactoryInstanceRegistry struct {
	*DefaultSingletonInstanceRegistry

	factoryInstanceCache map[string]interface{}
}

func NewFactoryInstanceRegistry() *FactoryInstanceRegistry {
	return &FactoryInstanceRegistry{
		DefaultSingletonInstanceRegistry: NewDefaultSingletonInstanceRegistry(),
		factoryInstanceCache:             make(map[string]interface{}),
	}
}

func (r *FactoryInstanceRegistry) getCachedObjectForFactoryInstance(name string) (interface{}, bool) {
	r.singletonMutex.RLock()
	obj, exists := r.factoryInstanceCache[name]
	r.singletonMutex.RUnlock()

	return obj, exists
}

func (r *FactoryInstanceRegistry) getObjectFromFactoryInstance(factory FactoryInstance, name string) (obj interface{}, err error) {
	if factory.IsSingleton() && r.ContainsSingleton(name) {
		var exists bool

		obj, exists = r.getCachedObjectForFactoryInstance(name)

		if !exists {
			obj, err = factory.GetInstance()

			if err == nil {
				r.singletonMutex.Lock()
				r.factoryInstanceCache[name] = obj
				r.singletonMutex.Unlock()
			}
		}
	} else {
		obj, err = factory.GetInstance()
	}

	return
}

func (r *FactoryInstanceRegistry) removeFactoryInstance(name string) {
	r.removeSingleton(name)

	r.singletonMutex.Lock()
	delete(r.factoryInstanceCache, name)
	r.singletonMutex.Unlock()
}

/*
type AbstractInstanceFactory struct {
	*FactoryInstanceRegistry
	InstanceFactoryUtils

	parentInstanceFactory InstanceFactory
}

var (
	_ = (ConfigurableInstanceFactory)((*AbstractInstanceFactory)(nil))
)

func NewAbstractInstanceFactory() *AbstractInstanceFactory {
	return &AbstractInstanceFactory{
		FactoryInstanceRegistry: NewFactoryInstanceRegistry(),
	}
}

func (f *AbstractInstanceFactory) ContainsInstance(name string) bool {
	name = f.transformedInstanceName(name)

	if f.ContainsSingleton(name) || f.ContainsInstanceDefinition(name) {
		return !f.isFactoryDereference(name) || isFactoryInstance(name)
	}
}

func (f *AbstractInstanceFactory) GetAliases(name string) []string {}

func (f *AbstractInstanceFactory) GetInstanceByType(t reflect.Type, args ...interface{}) interface{} {}

func (f *AbstractInstanceFactory) GetInstanceByName(name string, args ...interface{}) interface{} {}

func (f *AbstractInstanceFactory) GetType(name string) reflect.Type {}

func (f *AbstractInstanceFactory) IsPrototype(name string) bool {}

func (f *AbstractInstanceFactory) IsSingleton(name string) bool {}

func (f *AbstractInstanceFactory) IsFactoryInstance(name string) bool {
	name = f.transformedInstanceName(name)

	if obj, exists := f.getSingleton(name, false); exists {
		_, ok := obj.(FactoryInstance)

		return ok
	}

	if f.ContainsSingleton(name) {
		return false
	}

	if !f.ContainsBeanDefinition(name) {
		if factory, ok := f.parentInstanceFactory.(ConfigurableBeanFactory); ok {
			return factory.IsFactoryInstance(name)
		}
	}

	return false
}

func (f *AbstractInstanceFactory) transformedInstanceName(name string) string {
	return f.canonicalName(f.InstanceFactoryUtils.transformedInstanceName(name))
}

type DefaultListableInstanceFactory struct {
	*AbstractInstanceFactory
}

func NewDefaultListableInstanceFactory() *DefaultListableInstanceFactory {
	return &DefaultListableInstanceFactory{
		AbstractInstanceFactory: NewAbstractInstanceFactory(),
	}
}
*/
