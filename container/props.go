package container

import (
	"fmt"
	"strings"
)

type PropertySource interface {
	Name() string

	ContainsProperty(key string) bool

	GetProperty(key string) (string, bool)
}

// Interface for resolving properties against any underlying source.
type PropertyResolver interface {
	// Return whether the given property key is available for resolution, i.e.,
	// the value for the given key is not {@code null}.
	ContainsProperty(key string) bool

	// Return the property value associated with the given key
	GetProperty(key, defaultValue string) (string, bool)

	// Return the property value associated with the given key or panic
	GetRequiredProperty(key string) string
}

type ConfigurablePropertyResolver interface {
	PropertyResolver

	// Specify which properties must be present, to be verified by ValidateRequiredProperties().
	SetRequiredProperties(requiredProperties ...string)

	// Validate that each of the properties specified by SetRequiredProperties is present and resolves to a non-Nil value.
	ValidateRequiredProperties() error
}

type ErrMissingRequiredProperties []string

func (e ErrMissingRequiredProperties) Error() string {
	return fmt.Sprintf("The following properties were declared as required but could not be resolved: %s", strings.Join(e, ","))
}

type PropertySourcesPropertyResolver struct {
	propertySources    []PropertySource
	requiredProperties map[string]bool
}

var _ = (ConfigurablePropertyResolver)((*PropertySourcesPropertyResolver)(nil))

func NewPropertySourcesPropertyResolver() *PropertySourcesPropertyResolver {
	return &PropertySourcesPropertyResolver{
		requiredProperties: make(map[string]bool),
	}
}

func (r *PropertySourcesPropertyResolver) ContainsProperty(key string) bool {
	for _, propertySource := range r.propertySources {
		if propertySource.ContainsProperty(key) {
			return true
		}
	}

	return false
}

func (r *PropertySourcesPropertyResolver) GetProperty(key, defaultValue string) (string, bool) {
	for _, propertySource := range r.propertySources {
		if value, exists := propertySource.GetProperty(key); exists {
			return value, true
		}
	}

	return defaultValue, false
}

func (r *PropertySourcesPropertyResolver) GetRequiredProperty(key string) string {
	if value, exists := r.GetProperty(key, ""); exists {
		return value
	}

	panic(fmt.Errorf("required key [%s] not found", key))
}

func (r *PropertySourcesPropertyResolver) SetRequiredProperties(requiredProperties ...string) {
	for _, name := range requiredProperties {
		r.requiredProperties[name] = true
	}
}

func (r *PropertySourcesPropertyResolver) ValidateRequiredProperties() error {
	var names []string

	for name, required := range r.requiredProperties {
		if !required {
			continue
		}

		if value, exists := r.GetProperty(name, ""); !exists || value == "" {
			names = append(names, name)
		}
	}

	return ErrMissingRequiredProperties(names)
}
