package container

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SingletonRegistry interface {
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

type DefaultSingletonRegistry struct {
	lock                          sync.RWMutex
	singletonObjects              map[string]interface{}
	singletonFactories            map[string]ObjectFactory
	earlySingletonObjects         map[string]interface{}
	registeredSingletons          map[string]time.Time
	singletonsCurrentlyInCreation map[string]bool
	inCreationCheckExclusions     map[string]bool
}

func NewDefaultSingletonRegistry() *DefaultSingletonRegistry {
	return &DefaultSingletonRegistry{
		singletonObjects:              make(map[string]interface{}),
		singletonFactories:            make(map[string]ObjectFactory),
		earlySingletonObjects:         make(map[string]interface{}),
		registeredSingletons:          make(map[string]time.Time),
		singletonsCurrentlyInCreation: make(map[string]bool),
		inCreationCheckExclusions:     make(map[string]bool),
	}
}

func (r *DefaultSingletonRegistry) ContainsSingleton(name string) bool {
	r.lock.RLock()
	_, exists := r.singletonObjects[name]
	r.lock.RUnlock()

	return exists
}

func (r *DefaultSingletonRegistry) GetSingleton(name string) (interface{}, error) {
	return r.getSingleton(name, true)
}

func (r *DefaultSingletonRegistry) getSingleton(name string, allowEarlyReference bool) (interface{}, error) {
	r.lock.RLock()
	obj, exists := r.singletonObjects[name]
	r.lock.RUnlock()

	if !exists && r.isSingletonCurrentlyInCreation(name) {
		r.lock.RLock()
		obj, exists = r.earlySingletonObjects[name]
		r.lock.RUnlock()

		if !exists && allowEarlyReference {
			r.lock.RLock()
			factory, exists := r.singletonFactories[name]
			r.lock.RUnlock()

			if !exists || factory == nil {
				return nil, fmt.Errorf("singleton '%s' without instance or factory", name)
			}

			obj, err := factory.GetInstance()

			if err != nil {
				return nil, err
			}

			r.lock.Lock()
			r.earlySingletonObjects[name] = obj
			delete(r.singletonFactories, name)
			r.lock.Unlock()

			return obj, nil
		}
	}

	return obj, nil
}

func (r *DefaultSingletonRegistry) setCurrentlyInCreation(name string, inCreation bool) {
	if !inCreation {
		r.inCreationCheckExclusions[name] = true
	} else {
		delete(r.inCreationCheckExclusions, name)
	}
}

func (r *DefaultSingletonRegistry) isCurrentlyInCreation(name string) bool {
	_, exists := r.inCreationCheckExclusions[name]

	return !exists && r.isSingletonCurrentlyInCreation(name)
}

func (r *DefaultSingletonRegistry) isSingletonCurrentlyInCreation(name string) bool {
	_, exists := r.singletonsCurrentlyInCreation[name]

	return exists
}

func (r *DefaultSingletonRegistry) beforeSingletonCreation(name string) {
	_, exists := r.inCreationCheckExclusions[name]

	if !exists {
		r.singletonsCurrentlyInCreation[name] = true
	}
}

func (r *DefaultSingletonRegistry) afterSingletonCreation(name string) {
	_, exists := r.inCreationCheckExclusions[name]

	if !exists {
		delete(r.singletonsCurrentlyInCreation, name)
	}
}

func (r *DefaultSingletonRegistry) SingletonCount() int {
	r.lock.RLock()
	count := len(r.registeredSingletons)
	r.lock.RUnlock()

	return count
}

func (r *DefaultSingletonRegistry) SingletonNames() (names []string) {
	r.lock.RLock()

	for name, _ := range r.registeredSingletons {
		names = append(names, name)
	}

	r.lock.RUnlock()

	return
}

func (r *DefaultSingletonRegistry) RegisterSingleton(name string, obj interface{}) error {
	if name == "" {
		return errors.New("'name' must not be null")
	}

	if r.ContainsSingleton(name) {
		return fmt.Errorf("Could not register object under name '%s': there is already object bound", name)
	}

	r.addSingleton(name, obj)

	return nil
}

func (r *DefaultSingletonRegistry) addSingleton(name string, obj interface{}) {
	r.lock.Lock()
	r.singletonObjects[name] = obj
	delete(r.singletonFactories, name)
	r.registeredSingletons[name] = time.Now()
	r.lock.Unlock()
}

func (r *DefaultSingletonRegistry) addSingletonFactory(name string, factory ObjectFactory) error {
	if name == "" {
		return errors.New("'name' must not be null")
	}

	if factory == nil {
		return errors.New("'factory' must not be null")
	}

	if !r.ContainsSingleton(name) {
		r.lock.Lock()
		r.singletonFactories[name] = factory
		r.registeredSingletons[name] = time.Now()
		r.lock.Unlock()
	}

	return nil
}

func (r *DefaultSingletonRegistry) removeSingleton(name string) {
	r.lock.Lock()
	delete(r.singletonObjects, name)
	delete(r.singletonFactories, name)
	delete(r.registeredSingletons, name)
	r.lock.Unlock()
}
