package providers

import (
	"fmt"
	"os"
	"strings"
)

// EnvProvider loads values from the environment
type EnvProvider struct {
	EnvironSlice []string
}

var (
	_ Provider        = EnvProvider{}
	_ BuiltInProvider = EnvProvider{}
)

// Provide converts the key to upper case, replaces spaces and dashes with underscores and loads the value from the environment
func (p EnvProvider) Provide(key string) (string, error) {
	convertedKey := strings.ToUpper(key)
	convertedKey = strings.Replace(convertedKey, "-", "_", -1)
	convertedKey = strings.Replace(convertedKey, " ", "_", -1)

	if p.EnvironSlice != nil {
		return p.provideFromSlice(key, convertedKey, p.EnvironSlice)
	}

	return p.provideFromEnvironment(key, convertedKey)
}

func (p EnvProvider) provideFromEnvironment(key, convertedKey string) (string, error) {
	value, ok := os.LookupEnv(convertedKey)
	if !ok {
		return "", fmt.Errorf("unable to load value %s from environment variable %s", key, convertedKey)
	}

	return value, nil
}

func (p EnvProvider) provideFromSlice(key, convertedKey string, environ []string) (string, error) {
	for _, v := range environ {
		s := strings.SplitN(v, `=`, 2)
		if len(s) < 2 {
			continue
		}

		if s[0] == convertedKey {
			return s[1], nil
		}
	}

	return "", fmt.Errorf("unable to load value %v from environment variable %v", key, convertedKey)
}

// Tag returns the struct tag to look up environment variables
func (p EnvProvider) Tag() string {
	return "env"
}

// Precedence returns the precedence for the EnvProvider
func (p EnvProvider) Precedence() int {
	return 20
}

// IsBuiltInProvider implements a built-in provider interface
func (p EnvProvider) IsBuiltInProvider() bool {
	return true
}
