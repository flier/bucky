package container

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SingletonInstanceRegistry interface {
	// Check if this registry contains a singleton instance with the given name.
	ContainsSingleton(name string) bool

	// Return the (raw) singleton object registered under the given name.
	GetSingleton(name string) (interface{}, error)

	// Return the number of singleton instances registered in this registry.
	SingletonCount() int

	// Return the names of singleton instances registered in this registry.
	SingletonNames() []string

	// Register the given existing object as singleton in the registry, under the given bean name.
	RegisterSingleton(name string, obj interface{}) error
}

type DefaultSingletonInstanceRegistry struct {
	AliasRegistry

	singletonMutex                sync.RWMutex
	singletonObjects              map[string]interface{}
	singletonFactories            map[string]ObjectFactory
	earlySingletonObjects         map[string]interface{}
	registeredSingletons          map[string]time.Time
	singletonsCurrentlyInCreation map[string]bool
	inCreationCheckExclusions     map[string]bool
}

var (
	_ = (SingletonInstanceRegistry)((*DefaultSingletonInstanceRegistry)(nil))
)

func NewDefaultSingletonInstanceRegistry() *DefaultSingletonInstanceRegistry {
	return &DefaultSingletonInstanceRegistry{
		AliasRegistry:                 NewSimpleAliasRegistry(),
		singletonObjects:              make(map[string]interface{}),
		singletonFactories:            make(map[string]ObjectFactory),
		earlySingletonObjects:         make(map[string]interface{}),
		registeredSingletons:          make(map[string]time.Time),
		singletonsCurrentlyInCreation: make(map[string]bool),
		inCreationCheckExclusions:     make(map[string]bool),
	}
}

func (r *DefaultSingletonInstanceRegistry) ContainsSingleton(name string) bool {
	r.singletonMutex.RLock()
	_, exists := r.singletonObjects[name]
	r.singletonMutex.RUnlock()

	return exists
}

func (r *DefaultSingletonInstanceRegistry) GetSingleton(name string) (interface{}, error) {
	return r.getSingleton(name, true)
}

func (r *DefaultSingletonInstanceRegistry) getSingleton(name string, allowEarlyReference bool) (interface{}, error) {
	r.singletonMutex.RLock()
	obj, exists := r.singletonObjects[name]
	r.singletonMutex.RUnlock()

	if !exists && r.isSingletonCurrentlyInCreation(name) {
		r.singletonMutex.RLock()
		obj, exists = r.earlySingletonObjects[name]
		r.singletonMutex.RUnlock()

		if !exists && allowEarlyReference {
			r.singletonMutex.RLock()
			factory, exists := r.singletonFactories[name]
			r.singletonMutex.RUnlock()

			if !exists || factory == nil {
				return nil, fmt.Errorf("singleton '%s' without instance or factory", name)
			}

			obj, err := factory.GetInstance()

			if err != nil {
				return nil, err
			}

			r.singletonMutex.Lock()
			r.earlySingletonObjects[name] = obj
			delete(r.singletonFactories, name)
			r.singletonMutex.Unlock()

			return obj, nil
		}
	}

	return obj, nil
}

func (r *DefaultSingletonInstanceRegistry) setCurrentlyInCreation(name string, inCreation bool) {
	if !inCreation {
		r.inCreationCheckExclusions[name] = true
	} else {
		delete(r.inCreationCheckExclusions, name)
	}
}

func (r *DefaultSingletonInstanceRegistry) isCurrentlyInCreation(name string) bool {
	_, exists := r.inCreationCheckExclusions[name]

	return !exists && r.isSingletonCurrentlyInCreation(name)
}

func (r *DefaultSingletonInstanceRegistry) isSingletonCurrentlyInCreation(name string) bool {
	_, exists := r.singletonsCurrentlyInCreation[name]

	return exists
}

func (r *DefaultSingletonInstanceRegistry) beforeSingletonCreation(name string) {
	_, exists := r.inCreationCheckExclusions[name]

	if !exists {
		r.singletonsCurrentlyInCreation[name] = true
	}
}

func (r *DefaultSingletonInstanceRegistry) afterSingletonCreation(name string) {
	_, exists := r.inCreationCheckExclusions[name]

	if !exists {
		delete(r.singletonsCurrentlyInCreation, name)
	}
}

func (r *DefaultSingletonInstanceRegistry) SingletonCount() int {
	r.singletonMutex.RLock()
	count := len(r.registeredSingletons)
	r.singletonMutex.RUnlock()

	return count
}

func (r *DefaultSingletonInstanceRegistry) SingletonNames() (names []string) {
	r.singletonMutex.RLock()

	for name, _ := range r.registeredSingletons {
		names = append(names, name)
	}

	r.singletonMutex.RUnlock()

	return
}

func (r *DefaultSingletonInstanceRegistry) RegisterSingleton(name string, obj interface{}) error {
	if name == "" {
		return errors.New("'name' must not be null")
	}

	if r.ContainsSingleton(name) {
		return fmt.Errorf("Could not register object under name '%s': there is already object bound", name)
	}

	r.addSingleton(name, obj)

	return nil
}

func (r *DefaultSingletonInstanceRegistry) addSingleton(name string, obj interface{}) {
	r.singletonMutex.Lock()
	r.singletonObjects[name] = obj
	delete(r.singletonFactories, name)
	r.registeredSingletons[name] = time.Now()
	r.singletonMutex.Unlock()
}

func (r *DefaultSingletonInstanceRegistry) addSingletonFactory(name string, factory ObjectFactory) error {
	if name == "" {
		return errors.New("'name' must not be null")
	}

	if factory == nil {
		return errors.New("'factory' must not be null")
	}

	if !r.ContainsSingleton(name) {
		r.singletonMutex.Lock()
		r.singletonFactories[name] = factory
		r.registeredSingletons[name] = time.Now()
		r.singletonMutex.Unlock()
	}

	return nil
}

func (r *DefaultSingletonInstanceRegistry) removeSingleton(name string) {
	r.singletonMutex.Lock()
	delete(r.singletonObjects, name)
	delete(r.singletonFactories, name)
	delete(r.registeredSingletons, name)
	r.singletonMutex.Unlock()
}
