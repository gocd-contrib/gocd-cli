package materials

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/pflag"
)

type hash map[string]interface{}

func (h hash) SetRequiredString(key, val string) error {
	if `` == val {
		return fmt.Errorf(`%q is a required key`, val)
	}
	h[key] = val
	return nil
}

func (h hash) SetStringWithDefault(key, val, defaultVal string) {
	if `` == val {
		val = defaultVal
	}
	h[key] = val
}

func (h hash) SetStringIfFlagSet(key, val string, flag *pflag.Flag) {
	if flag.Changed {
		h[key] = val
	}
}

func (h hash) SetBoolIfFlagSet(key string, val bool, flag *pflag.Flag) {
	if flag.Changed {
		h[key] = val
	}
}

func (h hash) SetBool(key string, val bool) {
	h[key] = val
}

func (h hash) Equivalent(other hash) bool {
	if len(h) != len(other) {
		return false
	}

	for key, otherval := range other {
		if myval, ok := h[key]; !ok {
			return false
		} else {
			switch otherval.(type) {
			case string, bool, int, int64, float64:
				if myval != otherval {
					return false
				}
			case hash:
				if sub, ok := myval.(hash); !ok {
					return false
				} else {
					if !hash(sub).Equivalent(otherval.(hash)) {
						return false
					}
				}
			default:
				log.Fatalf("`hash` type does not know how to compare values of type %T; please add a case for this in hash.Equivalent()", otherval)
				return false
			}
		}
	}

	return true
}

// A strange concept in Golang, but some of the material attributes may be returned as `null`
// from GoCD even if they usually refer to strings. The `name` key comes to mind, and we need
// to allow this to be set.
func (h hash) copyNillableStrTo(dest hash, key string) error {
	if val, ok := h[key]; ok {
		if val == nil {
			dest[key] = nil
			return nil
		}

		if s, ok := val.(string); ok {
			dest[key] = s
		} else {
			return fmt.Errorf(`key %q does not hold a string; type: %T, value: %v`, key, val, val)
		}
	}
	return nil
}

func (h hash) copyStrIfPresentTo(dest hash, key string) error {
	if val, ok := h[key]; ok {
		if s, ok := val.(string); ok {
			dest[key] = s
		} else {
			return fmt.Errorf(`key %q does not hold a string; type: %T, value: %v`, key, val, val)
		}
	}
	return nil
}

func (h hash) copyStrOrDefaultTo(dest hash, key, defaultVal string) error {
	if val, ok := h[key]; ok {
		if s, ok := val.(string); ok {
			dest[key] = s
		} else {
			return fmt.Errorf(`key %q does not hold a string; type: %T, value: %v`, key, val, val)
		}
	} else {
		dest[key] = defaultVal
	}
	return nil
}

func (h hash) copyBoolIfPresentTo(dest hash, key string) error {
	if val, ok := h[key]; ok {
		if b, ok := val.(bool); ok {
			dest[key] = b
		} else {
			return fmt.Errorf(`key %q does not hold a bool; type: %T, value: %v`, key, val, val)
		}
	}
	return nil
}

func (h hash) copyBoolOrDefaultTo(dest hash, key string, defaultVal bool) error {
	if val, ok := h[key]; ok {
		if b, ok := val.(bool); ok {
			dest[key] = b
		} else {
			return fmt.Errorf(`key %q does not hold a bool; type: %T, value: %v`, key, val, val)
		}
	} else {
		dest[key] = defaultVal
	}
	return nil
}

func (h hash) subHash(key string) (hash, error) {
	if v, ok := h[key]; ok {
		if sub, ok := v.(map[string]interface{}); ok {
			return hash(sub), nil
		} else {
			return nil, fmt.Errorf(`value at key %q is not a map[string]interface{}, but %T %v`, key, v, v)
		}
	} else {
		return nil, fmt.Errorf(`key %q is not present`, key)
	}
}

func (h hash) string(key string) (string, error) {
	if v, ok := h[key]; ok {
		if str, ok := v.(string); ok {
			return str, nil
		} else {
			return ``, fmt.Errorf(`value at key %q is not a string, but %T %v`, key, v, v)
		}
	} else {
		return ``, fmt.Errorf(`key %q is not present`, key)
	}
}

func (h hash) String() string {
	b, _ := json.Marshal(h)
	return string(b)
}
