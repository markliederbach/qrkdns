package configr

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/markliederbach/configr/providers"
	"github.com/mitchellh/reflectwalk"
	"github.com/spf13/cast"
)

var (
	_ Walker = &DefaultWalker{}
)

// Walker defines methods necessary to integrate with reflectwalk
type Walker interface {
	reflectwalk.StructWalker

	// Err returns an error if any of the included providers had an error
	Err() error

	// Providers returns all providers for a walker
	Providers() map[string]providers.Provider
}

// DefaultWalker stores providers and results from loading the configuration
type DefaultWalker struct {
	errors    []error
	providers map[string]providers.Provider
}

// NewWalker returns a new confiration Walker with built-in providers
func NewWalker(opts []LoadOption) Walker {
	var walker Walker = &DefaultWalker{
		errors: []error{},
		providers: map[string]providers.Provider{
			"env":     providers.EnvProvider{},
			"file":    providers.FileProvider{},
			"default": providers.DefaultValueProvider{},
		},
	}
	for _, opt := range opts {
		walker = opt(walker)
	}

	return walker
}

// Err returns an error if any of the included providers had an error
func (w *DefaultWalker) Err() error {
	if len(w.errors) == 0 {
		return nil
	}

	errorStrings := []string{}
	for _, err := range w.errors {
		errorStrings = append(errorStrings, err.Error())
	}

	return fmt.Errorf("encountered %d errors loading config: %s", len(w.errors), strings.Join(errorStrings, ", "))
}

// Providers returns all providers for a walker
func (w *DefaultWalker) Providers() map[string]providers.Provider {
	return w.providers
}

// Struct implements reflectwalk.StructWalker interface
func (w *DefaultWalker) Struct(obj reflect.Value) error {
	return nil
}

// StructField implements reflectwalk.StructWalker interface
func (w *DefaultWalker) StructField(field reflect.StructField, value reflect.Value) error {
	err := w.populateField(field, value)
	if err != nil {
		w.errors = append(w.errors, fmt.Errorf("error loading config %s: %+v", field.Name, err))
	}
	return nil
}

func (w *DefaultWalker) populateField(field reflect.StructField, value reflect.Value) error {
	var err error
	// apply providers in precedence order, bailing out after the first one that succeeds
	for _, provider := range w.providersForField(field) {
		err = w.applyProvider(provider, field, value)

		if err == nil {
			return nil
		}
	}

	return err
}

func (w *DefaultWalker) providersForField(field reflect.StructField) []providers.Provider {
	providers := []providers.Provider{}

	// first determine provides that apply to this field
	for _, provider := range w.providers {
		if _, ok := field.Tag.Lookup(provider.Tag()); ok {
			providers = append(providers, provider)
		}
	}

	// sort them by precedence: higher precedence gets applied first
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Precedence() > providers[j].Precedence()
	})

	return providers
}

func (w *DefaultWalker) applyProvider(provider providers.Provider, field reflect.StructField, value reflect.Value) error {
	key := field.Tag.Get(provider.Tag())
	newValue, err := provider.Provide(key)
	if err != nil {
		return err
	}

	err = w.setValue(value, newValue)
	if err != nil {
		return err
	}

	return nil
}

func (w *DefaultWalker) setValue(value reflect.Value, newValue string) error {
	switch t := value.Interface().(type) {
	case string:
		value.SetString(newValue)

	case []byte:
		value.SetBytes([]byte(newValue))

	case []string:
		parts := strings.Split(newValue, ",")
		result := reflect.MakeSlice(reflect.TypeOf(parts), 0, len(parts))
		for _, part := range parts {
			result = reflect.Append(result, reflect.ValueOf(part))
		}
		value.Set(result)

	case int, int8, int16, int32, int64:
		intVal, err := cast.ToInt64E(newValue)
		if err != nil {
			return fmt.Errorf("unable to convert %s to int", newValue)
		}
		value.SetInt(intVal)

	case uint, uint8, uint16, uint32, uint64:
		uintVal, err := cast.ToUint64E(newValue)
		if err != nil {
			return fmt.Errorf("unable to convert %s to uint", newValue)
		}
		value.SetUint(uintVal)

	case float32, float64:
		floatVal, err := cast.ToFloat64E(newValue)
		if err != nil {
			return fmt.Errorf("unable to convert %s to float", newValue)
		}
		value.SetFloat(floatVal)

	case time.Duration:
		durationVal, err := cast.ToDurationE(newValue)
		if err != nil {
			return fmt.Errorf("unable to convert %s to duration", newValue)
		}
		value.Set(reflect.ValueOf(durationVal))

	case bool:
		boolVal, err := cast.ToBoolE(newValue)
		if err != nil {
			return fmt.Errorf("unable to convert %s to bool", newValue)
		}
		value.SetBool(boolVal)

	default:
		return fmt.Errorf("unsupported type: %+v", t)
	}

	return nil
}
