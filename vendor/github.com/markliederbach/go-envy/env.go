package envy

import (
	"errors"
	"os"
)

// MockEnv allows a test to temporarily set environment variables
type MockEnv struct {
	isLoaded  bool
	variables map[string]string
	existing  map[string]string
}

// Load sets environment variables from a given map
func (m *MockEnv) Load(variables map[string]string) error {
	if m.isLoaded {
		return errors.New("mock environment is already loaded")
	}

	m.existing = make(map[string]string)
	m.variables = variables

	for key, value := range m.variables {
		currentValue, _ := os.LookupEnv(key)
		m.existing[key] = currentValue

		if value != "" {
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		} else {
			_ = os.Unsetenv(key)
		}
	}
	m.isLoaded = true
	return nil
}

// Restore resets the environment variables back to their value
// prior to when Load was called
func (m *MockEnv) Restore() {
	if !m.isLoaded {
		panic(errors.New("environment not loaded"))
	}

	for key, value := range m.existing {
		if value != "" {
			err := os.Setenv(key, value)
			if err != nil {
				panic(err)
			}
		} else {
			_ = os.Unsetenv(key)
		}
	}
	m.existing = make(map[string]string)
	m.variables = make(map[string]string)
}
