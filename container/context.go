package container

import (
	"time"
)

type ApplicationContext interface {
	EnvironmentCapable
	ListableInstanceFactory
	HierarchicalInstanceFactory

	// Return a name for the deployed application that this context belongs to.
	ApplicationName() string

	// Return a friendly name for this context.
	DisplayName() string

	// Return the unique id of this application context.
	Id() string

	// Return the parent context, or nil if there is no parent and this is the root of the context hierarchy.
	Parent() ApplicationContext

	// Return the timestamp when this context was first loaded.
	StartupTime() time.Time

	// Expose AutowireCapableBeanFactory functionality for this context.
	AutowireCapableInstanceFactory() AutowireCapableInstanceFactory
}

type ConfigurableApplicationContext interface {
}

type AbstractApplicationContext struct {
	ConfigurableListableInstanceFactory

	appName         string
	id              string
	displayName     string
	startupTime     time.Time
	parent          ApplicationContext
	env             ConfigurableEnvironment
	instanceFactory ConfigurableListableInstanceFactory
}

var (
	_ = (ConfigurableApplicationContext)((*AbstractApplicationContext)(nil))
	_ = (ApplicationContext)((*AbstractApplicationContext)(nil))
)

func NewAbstractApplicationContext() *AbstractApplicationContext {
	return &AbstractApplicationContext{}
}

func (c *AbstractApplicationContext) Environment() Environment { return c.env }

func (c *AbstractApplicationContext) ApplicationName() string { return c.appName }

func (c *AbstractApplicationContext) DisplayName() string { return c.displayName }

func (c *AbstractApplicationContext) Id() string { return c.id }

func (c *AbstractApplicationContext) Parent() ApplicationContext { return c.parent }

func (c *AbstractApplicationContext) StartupTime() time.Time { return c.startupTime }

func (c *AbstractApplicationContext) AutowireCapableInstanceFactory() AutowireCapableInstanceFactory {
	return c.instanceFactory
}

type GenericApplicationContext struct {
	*AbstractApplicationContext
}

func NewGenericApplicationContext() *GenericApplicationContext {
	return &GenericApplicationContext{
		AbstractApplicationContext: NewAbstractApplicationContext(),
	}
}
