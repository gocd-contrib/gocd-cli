package configrepo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
)

type PropertySet map[string]Property

func (k PropertySet) Add(p Property) {
	if _, ok := k[p.Name()]; ok {
		utils.Errfln(`property %q declared multiple times; overriding...`, p.Name())
	}
	k[p.Name()] = p
}

func (k PropertySet) Append(v ...Property) PropertySet {
	for _, p := range v {
		k.Add(p)
	}
	return k
}

func (k PropertySet) Value() []Property {
	vals := make([]Property, 0)
	for _, val := range k {
		vals = append(vals, val)
	}
	return vals
}

func NewPropertySetValue(propertyFactory func(string, string) Property) *PropertySetFlag {
	kva := make(PropertySet)
	return MakePropertySetValue(&kva, propertyFactory)
}

func MakePropertySetValue(initial *PropertySet, propertyFactory func(string, string) Property) *PropertySetFlag {
	return &PropertySetFlag{value: initial, propertyFactory: propertyFactory, changed: false}
}

type PropertySetFlag struct {
	propertyFactory func(string, string) Property
	value           *PropertySet
	changed         bool
}

func (k *PropertySetFlag) convert(val string) (Property, error) {
	if strings.Count(val, `:`) == 0 {
		return nil, fmt.Errorf(`%s must be formatted as key:value`, val)
	}

	// split only on first `:`; this will not allow any escaping of `:`,
	// so keys should NOT have colon characters. will not address escaping
	// until we find a need for it to avoid complexity.
	pair := strings.SplitN(val, `:`, 2)
	return k.propertyFactory(pair[0], pair[1]), nil
}

func (k *PropertySetFlag) Value() PropertySet {
	return *k.value
}

func (k *PropertySetFlag) Set(val string) error {
	v, err := k.convert(val)
	if err != nil {
		return err
	}

	(*k.value).Add(v)
	k.changed = true
	return nil
}

func (k *PropertySetFlag) Type() string {
	return "key:value"
}

func (k *PropertySetFlag) String() string {
	if len(*k.value) == 0 {
		return `` // hack to prevent `(default [])` from appearing in the flag's usage message
	}
	return fmt.Sprintf(`%v`, (*k.value).Value())
}

type jsonFlag bool

func (f *jsonFlag) Type() string {
	return `bool`
}

func (f *jsonFlag) String() string {
	return strconv.FormatBool(bool(*f))
}

func (f *jsonFlag) Set(value string) error {
	b, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	*f = jsonFlag(b)

	if b {
		RootCmd.PersistentFlags().Set(`plugin-id`, `json.config.plugin`)
	}

	return nil
}

func (f *jsonFlag) IsBoolFlag() bool {
	return true
}

type yamlFlag bool

func (f *yamlFlag) Type() string {
	return `bool`
}

func (f *yamlFlag) String() string {
	return strconv.FormatBool(bool(*f))
}

func (f *yamlFlag) Set(value string) error {
	b, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	*f = yamlFlag(b)

	if b {
		RootCmd.PersistentFlags().Set(`plugin-id`, `yaml.config.plugin`)
	}

	return nil
}

func (f *yamlFlag) IsBoolFlag() bool {
	return true
}

type groovyFlag bool

func (f *groovyFlag) Type() string {
	return `bool`
}

func (f *groovyFlag) String() string {
	return strconv.FormatBool(bool(*f))
}

func (f *groovyFlag) Set(value string) error {
	b, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	*f = groovyFlag(b)

	if b {
		RootCmd.PersistentFlags().Set(`plugin-id`, `cd.go.contrib.plugins.configrepo.groovy`)
	}

	return nil
}

func (f *groovyFlag) IsBoolFlag() bool {
	return true
}

func newJsonFlag(b bool) *jsonFlag {
	return (*jsonFlag)(&b)
}

func newYamlFlag(b bool) *yamlFlag {
	return (*yamlFlag)(&b)
}

func newGroovyFlag(b bool) *groovyFlag {
	return (*groovyFlag)(&b)
}
