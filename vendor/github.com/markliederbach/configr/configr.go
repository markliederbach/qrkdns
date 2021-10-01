package configr

import (
	"fmt"

	"github.com/markliederbach/configr/providers"
	"github.com/mitchellh/reflectwalk"
)

// LoadOption is a function used to override the default
// configr loading behavior
type LoadOption func(Walker) Walker

// Load fills the provided data object based on various providers
func Load(data interface{}, opts ...LoadOption) error {
	configrWalker := NewWalker(opts)
	err := reflectwalk.Walk(data, configrWalker)
	if err != nil {
		return fmt.Errorf("error loading configuration: %+v", err)
	}
	return configrWalker.Err()
}

// WithProvider overrides a provider for a given struct tag
func WithProvider(provider providers.Provider) LoadOption {
	return func(w Walker) Walker {
		w.Providers()[provider.Tag()] = provider
		return w
	}
}

// WithoutBuiltInProviders will skip out the bit where configr loads the built
// in env and file providers, and instead will only load providers explicitly
// specified with WithProvider
func WithoutBuiltInProviders() LoadOption {
	return func(w Walker) Walker {
		for providerName, provider := range w.Providers() {
			if biProvider, ok := provider.(providers.BuiltInProvider); ok {
				if biProvider.IsBuiltInProvider() {
					delete(w.Providers(), providerName)
				}
			}
		}
		return w
	}
}
