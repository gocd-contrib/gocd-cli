package configrepo

import "strconv"

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

func newJsonFlag(b bool) *jsonFlag {
	return (*jsonFlag)(&b)
}

func newYamlFlag(b bool) *yamlFlag {
	return (*yamlFlag)(&b)
}
