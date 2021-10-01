package providers

// Provider defines the functions used to load a specific
// type of configuration.
type Provider interface {
	// Provide loads the value for a given key
	Provide(key string) (string, error)

	// Tag is the struct tag to look for configuration path (eg. name of environment variable)
	Tag() string

	// Precedence determines the order in which config is applied when multiple tags are present
	Precedence() int
}

// BuiltInProvider defines the necessary interface for all built-in providers
type BuiltInProvider interface {
	IsBuiltInProvider() bool
}
