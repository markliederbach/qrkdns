package providers

import (
	"io/ioutil"
)

// FileProvider loads configuration from a file on disk (relative or absolute path)
type FileProvider struct{}

var (
	_ Provider        = FileProvider{}
	_ BuiltInProvider = FileProvider{}
)

// Provide reads the file off disk, if it exists, and returns the contents cast as a string
func (p FileProvider) Provide(key string) (string, error) {
	fileBytes, err := ioutil.ReadFile(key)
	if err != nil {
		return "", err
	}
	return string(fileBytes), err
}

// Tag returns the struct tag to look up environment variables
func (p FileProvider) Tag() string {
	return "file"
}

// Precedence returns the precedence for the FileProvider
// More than default, but less than environment
func (p FileProvider) Precedence() int {
	return 15
}

// IsBuiltInProvider implements a built-in provider interface
func (p FileProvider) IsBuiltInProvider() bool {
	return true
}
