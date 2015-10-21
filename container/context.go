package container

import (
	"time"
)

type ApplicationContext interface {
	InstanceFactory

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
}

type AbstractApplicationContext struct {
	ApplicationContext
}

type GenericApplicationContext struct {
	AbstractApplicationContext
}
