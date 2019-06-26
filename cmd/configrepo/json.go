package configrepo

import (
	"encoding/json"
	"fmt"

	"github.com/gocd-contrib/gocd-cli/materials"
)

type ConfigRepo struct {
	Id            string             `json:"id"`
	PluginId      string             `json:"plugin_id"`
	Material      materials.Material `json:"material"`
	Configuration []Property         `json:"configuration"`
}

func NewConfigRepo(id, pluginId string, mat materials.Material, properties ...Property) *ConfigRepo {
	if properties == nil {
		properties = make([]Property, 0)
	}
	return &ConfigRepo{Id: id, PluginId: pluginId, Material: mat, Configuration: properties}
}

func (cr *ConfigRepo) UnmarshalJSON(data []byte) (err error) {
	var _m dict
	if err = json.Unmarshal(data, &_m); err != nil {
		return err
	}

	if cr.Id, err = _m.stringOrError(`id`); err != nil {
		return err
	}

	if cr.PluginId, err = _m.stringOrError(`plugin_id`); err != nil {
		return err
	}

	if m, err := _m.dictOrError(`material`); err != nil {
		return err
	} else {
		if cr.Material, err = materials.FromMap(m); err != nil {
			return err
		}
	}

	if _, ok := _m[`configuration`]; ok {
		if sl, err := _m.sliceOrError(`configuration`); err != nil {
			return err
		} else {
			for _, el := range sl {
				if d, ok := el.(map[string]interface{}); !ok {
					return fmt.Errorf(`configuration item is not a dict, but %T => %v`, el, el)
				} else {
					if prop, err := propertyFromMap(d); err == nil {
						cr.Configuration = append(cr.Configuration, prop)
					} else {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (cr *ConfigRepo) Equivalent(other *ConfigRepo) bool {
	if other == nil || cr.Id != other.Id {
		return false
	}

	if (cr.Material == nil && other.Material != nil) || (other.Material == nil && cr.Material != nil) {
		return false
	}

	if !cr.Material.Equivalent(other.Material) {
		return false
	}

	return true
}

func propertyFromMap(data map[string]interface{}) (prop Property, err error) {
	el := dict(data)
	if k, ok := el[`key`]; !ok {
		return nil, fmt.Errorf(`Configuration property is missing "key": %v`, el)
	} else {
		if _, ok := el[`encrypted_value`]; ok {
			if _, ok := el[`value`]; ok {
				return nil, fmt.Errorf(`Configuration property %q cannot have both "value" and "encrypted_value": %v`, k, el)
			}

			val := &SecretProperty{}
			if val.Key, err = el.stringOrError(`key`); err != nil {
				return nil, err
			}
			if val.EncryptedValue, err = el.stringOrError(`encrypted_value`); err != nil {
				return nil, err
			}
			return val, nil
		} else if _, ok := el[`value`]; ok {
			val := &PlainTextProperty{}
			if val.Key, err = el.stringOrError(`key`); err != nil {
				return nil, err
			}
			if val.Value, err = el.stringOrError(`value`); err != nil {
				return nil, err
			}
			return val, nil
		} else {
			return nil, fmt.Errorf(`Configuration property %q is missing a value: %v`, k, el)
		}
	}
}

type Property interface {
	Name() string
	String() string
	Equivalent(Property) bool
}

func NewPlainTextProperty(key, val string) Property {
	return &PlainTextProperty{Key: key, Value: val}
}

func NewSecretProperty(key, val string) Property {
	return &SecretProperty{Key: key, EncryptedValue: val}
}

type PlainTextProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p *PlainTextProperty) Name() string {
	return p.Key
}

func (p *PlainTextProperty) String() string {
	return p.Key + `: ` + p.Value
}

func (p *PlainTextProperty) Equivalent(other Property) bool {
	that, ok := other.(*PlainTextProperty)
	return ok && *p == *that
}

type SecretProperty struct {
	Key            string `json:"key"`
	EncryptedValue string `json:"encrypted_value"`
}

func (p *SecretProperty) Name() string {
	return p.Key
}

func (p *SecretProperty) String() string {
	return p.Key + `: ************`
}

func (p *SecretProperty) Equivalent(other Property) bool {
	that, ok := other.(*SecretProperty)
	return ok && *p == *that
}

type dict map[string]interface{}

func (d dict) sliceOrError(key string) ([]interface{}, error) {
	if v, ok := d[key]; ok {
		if sl, ok := v.([]interface{}); ok {
			return sl, nil
		} else {
			return nil, fmt.Errorf(`Value for key %q is not a slice, but a %T => %v`, key, v, v)
		}
	} else {
		return nil, fmt.Errorf(`dict is missing key %q`, key)
	}
}

func (d dict) dictOrError(key string) (dict, error) {
	if v, ok := d[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return dict(m), nil
		} else {
			return nil, fmt.Errorf(`Value for key %q is not a dict, but a %T => %v`, key, v, v)
		}
	} else {
		return nil, fmt.Errorf(`dict is missing key %q`, key)
	}
}

func (d dict) stringOrError(key string) (string, error) {
	if v, ok := d[key]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		} else {
			return "", fmt.Errorf(`Value for key %q is not a string, but a %T => %v`, key, v, v)
		}
	} else {
		return "", fmt.Errorf(`dict is missing key %q`, key)
	}
}
