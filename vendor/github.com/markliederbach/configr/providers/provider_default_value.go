package providers

// DefaultValueProvider provides a fallback for a given field
type DefaultValueProvider struct{}

var (
	_ Provider        = DefaultValueProvider{}
	_ BuiltInProvider = DefaultValueProvider{}
)

// Provide returns the value passed in
func (DefaultValueProvider) Provide(key string) (string, error) {
	return key, nil
}

// Tag is the struct tag where default values are defined
func (DefaultValueProvider) Tag() string {
	return "default"
}

// Precedence determines the order in which config is applied when multiple tags are present
func (DefaultValueProvider) Precedence() int {
	return 10
}

// IsBuiltInProvider implements a built-in provider interface
func (DefaultValueProvider) IsBuiltInProvider() bool {
	return true
}
