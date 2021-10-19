package envy

import "errors"

var (
	// ObjectChannels stores various objects by function name
	ObjectChannels map[string]chan interface{} = make(map[string]chan interface{})

	// ErrorChannels stores various errors by function name
	ErrorChannels map[string]chan error = make(map[string]chan error)

	// DefaultObjects stores various default objects by function name
	DefaultObjects map[string]interface{} = make(map[string]interface{})

	// DefaultErrors stores various default errors by function name
	DefaultErrors map[string]error = make(map[string]error)
)

// GetObject looks up the next object return from a channel, defaults if none exists
func GetObject(functionName string) interface{} {
	select {
	case obj := <-ObjectChannels[functionName]:
		return obj
	default:
		return DefaultObjects[functionName]
	}
}

// GetError looks up the next error return from a channel, defaults if none exists
func GetError(functionName string) error {
	select {
	case err := <-ErrorChannels[functionName]:
		return err
	default:
		return DefaultErrors[functionName]
	}
}

// AddObjectReturns allows a test to add one or more object returns for a function
func AddObjectReturns(functionName string, objs ...interface{}) error {
	for _, obj := range objs {
		if ch, ok := ObjectChannels[functionName]; !ok {
			return errors.New("Function channel does not exist")
		} else {
			ch <- obj
		}
	}
	return nil
}

// AddErrorReturns allows a test to add one or more error returns for a function
func AddErrorReturns(functionName string, errs ...error) error {
	for _, err := range errs {
		if ch, ok := ErrorChannels[functionName]; !ok {
			return errors.New("Function channel does not exist")
		} else {
			ch <- err
		}
	}
	return nil
}
