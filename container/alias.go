package container

import (
	"errors"
	"fmt"
	"sync"
)

// Common interface for managing aliases.
type AliasRegistry interface {
	// Given a name, register an alias for it.
	RegisterAlias(name, alias string) error

	// Remove the specified alias from this registry.
	RemoveAlias(name string) error

	// Determine whether this given name is defines as an alias
	IsAlias(name string) bool

	// Return the aliases for the given name, if defined.
	GetAliases(name string) []string
}

type SimpleAliasRegistry struct {
	lock                 sync.RWMutex
	aliases              map[string]string
	AllowAliasOverriding bool
}

func NewSimpleAliasRegistry() *SimpleAliasRegistry {
	return &SimpleAliasRegistry{
		aliases:              make(map[string]string),
		AllowAliasOverriding: true,
	}
}

func (r *SimpleAliasRegistry) RegisterAlias(name, alias string) error {
	if name == "" {
		return errors.New("'name' must not be empty")
	}

	if alias == "" {
		return errors.New("'alias' must not be empty")
	}

	if name == alias {
		r.lock.Lock()
		delete(r.aliases, name)
		r.lock.Unlock()
		return nil
	}

	if n, exists := r.aliases[alias]; exists {
		if n == name {
			// An existing alias - no need to re-register
			return nil
		}

		if !r.AllowAliasOverriding {
			return fmt.Errorf("Cannot register alias '%s' for name '%s', already registered for name `%s`", alias, name, n)
		}
	}

	if err := r.checkForAliasCircle(name, alias); err != nil {
		return err
	}

	r.lock.Lock()
	r.aliases[alias] = name
	r.lock.Unlock()

	return nil
}

func (r *SimpleAliasRegistry) checkForAliasCircle(name, alias string) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.hasAlias(alias, name) {
		return fmt.Errorf("Cannot register alias '%s' for name '%s', Circular reference - '%s' is a direct or indirect alias for '%s'", alias, name, name, alias)
	}

	return nil
}

func (r *SimpleAliasRegistry) hasAlias(name, alias string) bool {
	for key, value := range r.aliases {
		if value == name {
			return key == alias || r.hasAlias(key, alias)
		}
	}

	return false
}

func (r *SimpleAliasRegistry) RemoveAlias(alias string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, exists := r.aliases[alias]; !exists {
		return fmt.Errorf("No alias '%s' registered", alias)
	}

	delete(r.aliases, alias)

	return nil
}

func (r *SimpleAliasRegistry) IsAlias(alias string) bool {
	r.lock.RLock()
	_, exists := r.aliases[alias]
	r.lock.RUnlock()

	return exists
}

func (r *SimpleAliasRegistry) GetAliases(name string) []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.retrieveAliases(name)
}

func (r *SimpleAliasRegistry) retrieveAliases(name string) (result []string) {
	for alias, value := range r.aliases {
		if value == name {
			result = append(result, alias)
			result = append(result, r.retrieveAliases(alias)...)
		}
	}

	return
}
