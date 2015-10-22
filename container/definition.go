package container

type InstanceScope int

const (
	ScopeUnknown InstanceScope = iota
	ScopeSingleton
	ScopePrototype
)

// A InstanceDefinition describes a instance,
// which has property values, constructor argument values,
// and further information supplied by concrete implementations.
type InstanceDefinition interface {
	// Return the instance names that this instance depends on.
	DependsOn() []string

	// Return a human-readable description of this instance definition.
	Description() string

	// Return the factory instance name, if any.
	FactoryInstanceName() string

	// Return a factory method, if any.
	FactoryMethodName() string

	// Return the name of the parent definition of this instance definition, if any.
	ParentName() string

	// Return the name of the current target scope for this instance.
	Scope() InstanceScope

	// Return whether this bean should be lazily initialized, i.e. not eagerly instantiated on startup.
	IsLazyInit() bool

	// Return whether this instance is "abstract", that is, not meant to be instantiated.
	IsAbstract() bool

	// Return whether this a Prototype, with an independent instance returned for each call.
	IsPrototype() bool

	// Return whether this a Singleton, with a single, shared instance returned on all calls.
	IsSingleton() bool
}

type AbstractInstanceDefinition struct {
	depends         []string
	desc            string
	factoryInstance string
	factoryMethod   string
	parent          string
	scope           InstanceScope
	abstract        bool
	lazy            bool
}

var _ = (InstanceDefinition)((*AbstractInstanceDefinition)(nil))

func (d *AbstractInstanceDefinition) DependsOn() []string { return d.depends }

func (d *AbstractInstanceDefinition) Description() string { return d.desc }

func (d *AbstractInstanceDefinition) FactoryInstanceName() string { return d.factoryInstance }

func (d *AbstractInstanceDefinition) FactoryMethodName() string { return d.factoryMethod }

func (d *AbstractInstanceDefinition) ParentName() string { return d.parent }

func (d *AbstractInstanceDefinition) Scope() InstanceScope { return d.scope }

func (d *AbstractInstanceDefinition) IsAbstract() bool { return d.abstract }

func (d *AbstractInstanceDefinition) IsLazyInit() bool { return d.lazy }

func (d *AbstractInstanceDefinition) IsPrototype() bool { return d.scope == ScopePrototype }

func (d *AbstractInstanceDefinition) IsSingleton() bool { return d.scope == ScopeSingleton }

type RootBeanDefinition struct {
	*AbstractInstanceDefinition
}

func NewRootBeanDefinition() *RootBeanDefinition {
	return &RootBeanDefinition{
		AbstractInstanceDefinition: &AbstractInstanceDefinition{},
	}
}

type GenericInstanceDefinition struct {
	AbstractInstanceDefinition
}
