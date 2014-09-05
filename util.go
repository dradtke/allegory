package allegory

import (
	"errors"
	"fmt"
)

// errorize() takes any value (usually from recover()) and ensures
// that it either is an error value, or turns it into one.
func errorize(value interface{}) error {
	switch v := value.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return fmt.Errorf("%v", v)
	}
}
