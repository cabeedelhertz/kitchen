package service

import "context"

// PreStartHook is a function that will be called when the service is starting,
// but before it is receiving traffic
type PreStartHook func(context.Context) error

// PreShutdownHook is a function that will be called right before a service is
// stopped
type PreShutdownHook func() error

// ShutdownHook is a function that will be called when the service is stopping
type ShutdownHook func() error
