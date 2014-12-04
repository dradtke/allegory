package allegory

import (
	"bytes"
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"reflect"
	"strconv"
	"sync"
	"unicode"
)

// Returns true if the key is being held down.
func KeyDown(keyCode allegro.KeyCode) bool {
    return _pressedKeys[keyCode]
}

// After() takes a list of functions and kicks each one off in its own goroutine,
// then calls the callback once they've all finished. Everything is run
// in a separate goroutine, so After() returns almost immediately.
func After(routines []func(), callback func()) {
	var wg sync.WaitGroup
	wg.Add(len(routines))
	for _, routine := range routines {
		go func(f func()) {
			f()
			wg.Done()
		}(routine)
	}
	go func() {
		wg.Wait()
		callback()
	}()
}

// ReadConfig() reads the config section and saves all values it finds into
// dest. A "cfg" tag can specify which config value should be saved in
// that field. By default it will look for a field translated from snake-case
// to camel-case, e.g. hero_speed -> HeroSpeed.
func ReadConfig(cfg *allegro.Config, section string, dest interface{}) {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		panic("ReadConfig's `dest` must be a pointer!")
	}

	destVal = destVal.Elem()
	if destVal.Kind() != reflect.Struct {
		panic("ReadConfig's `dest` must point to a struct!")
	}

	n := destVal.NumField()
	for i := 0; i < n; i++ {
		field := destVal.Type().Field(i)
		fieldVal := destVal.Field(i)

		name := field.Tag.Get("cfg")
		if name == "" {
			name = field.Name
		}
		Debug(name)

		if val, err := cfg.Value(section, name); err == nil {
			saveToField(fieldVal, val)
			continue
		}

		name = camelToSnake(name)
		if val, err := cfg.Value(section, name); err == nil {
			err := saveToField(fieldVal, val)
			if err != nil {
				Fatal(err)
			}
			continue
		}
	}
}

// saveToField() saves `data` to `fieldVal`, converting it if necessary.
func saveToField(fieldVal reflect.Value, data string) error {
	switch fieldVal.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(data, 0, 64)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(int(i)))

	case reflect.Int8:
		i, err := strconv.ParseInt(data, 0, 8)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(int8(i)))

	case reflect.Int16:
		i, err := strconv.ParseInt(data, 0, 16)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(int16(i)))

	case reflect.Int32:
		i, err := strconv.ParseInt(data, 0, 32)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(int32(i)))

	case reflect.Int64:
		i, err := strconv.ParseInt(data, 0, 64)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(int64(i)))

	case reflect.Uint:
		i, err := strconv.ParseUint(data, 0, 64)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(uint(i)))

	case reflect.Uint8:
		i, err := strconv.ParseUint(data, 0, 8)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(uint8(i)))

	case reflect.Uint16:
		i, err := strconv.ParseUint(data, 0, 16)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(uint16(i)))

	case reflect.Uint32:
		i, err := strconv.ParseUint(data, 0, 32)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(uint32(i)))

	case reflect.Uint64:
		i, err := strconv.ParseUint(data, 0, 64)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(uint64(i)))

	case reflect.Float32:
		i, err := strconv.ParseFloat(data, 32)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(float32(i)))

	case reflect.Float64:
		i, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(float64(i)))

	default:
		return fmt.Errorf("tried to save config value to unsupported variable type: %s", fieldVal.Type().Name())
	}

	return nil
}

// camelToSnake() converts a camel-cased name to snake case.
func camelToSnake(str string) string {
	var snakeName bytes.Buffer
	for i, char := range str {
		if unicode.IsUpper(char) {
			if i > 0 {
				snakeName.WriteRune('_')
			}
			snakeName.WriteRune(unicode.ToLower(char))
		} else {
			snakeName.WriteRune(char)
		}
	}
	return snakeName.String()
}
