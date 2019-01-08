package cfg

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

func writeContent(fs afero.Fs, path, content string) error {
	if f, err := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.FileMode(0644)); err != nil {
		return err
	} else {
		defer f.Close()

		_, err = f.WriteString(content)
		return err
	}
}

func makeConf(content string) (*Config, error) {
	fs := afero.NewMemMapFs()
	v := newViper(fs)

	if err := writeContent(fs, TEST_CONF_FILE, content); err != nil {
		return nil, err
	} else {
		v.SetConfigFile(TEST_CONF_FILE)

		if err = v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{native: v, fs: fs}, nil
}

// creates a test Config instance, backed by a memory mapped afero.Fs
// if create == true, creates an empty config file on the afero.Fs
func testConf(create bool) *Config {
	fs := afero.NewMemMapFs()
	v := newViper(fs)

	v.AutomaticEnv()
	v.SetConfigName(CONFIG_FILENAME)
	v.SetConfigFile(TEST_CONF_FILE)

	if create {
		v.WriteConfigAs(TEST_CONF_FILE)
	}

	return &Config{native: v, fs: fs}
}

func serialize(t *testing.T, v interface{}) string {
	if b, err := yaml.Marshal(v); err != nil {
		t.Errorf("Error while trying to marshal %v: %v", v, err)
		t.FailNow()
		return ""
	} else {
		return string(b)
	}
}

type asserter struct {
	t *testing.T
}

func (a *asserter) deepEq(expected, actual interface{}) {
	a.t.Helper()
	a.eq(serialize(a.t, expected), serialize(a.t, actual))
}

func (a *asserter) configEq(expected dict, fs afero.Fs) {
	a.t.Helper()
	s, e := afero.ReadFile(fs, TEST_CONF_FILE)
	a.ok(e)

	exp, e1 := yaml.Marshal(expected)
	a.ok(e1)

	a.eq(string(exp), string(s))
}

func (a *asserter) eq(expected, actual interface{}) {
	a.t.Helper()
	if expected != actual {
		a.t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func (a *asserter) neq(expected, actual interface{}) {
	a.t.Helper()
	if expected == actual {
		a.t.Errorf("Expected %v to not equal %v", actual, expected)
	}
}

func (a *asserter) err(expected string, e error) {
	a.t.Helper()
	if nil == e {
		a.t.Errorf("Expected error %q, but got nil", expected)
		return
	}

	if e.Error() != expected {
		a.t.Errorf("Expected error %q, but got %q", expected, e)
	}
}

func (a *asserter) ok(err error) {
	a.t.Helper()
	if nil != err {
		a.t.Errorf("Expected no error, but got %v", err)
	}
}

func (a *asserter) is(b bool) {
	a.t.Helper()
	if !b {
		a.t.Errorf("Expected to be true")
	}
}

func (a *asserter) not(b bool) {
	a.t.Helper()
	if b {
		a.t.Errorf("Expected to be false")
	}
}

func asserts(t *testing.T) *asserter {
	return &asserter{t: t}
}
