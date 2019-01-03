package cfg

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

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

func (a *asserter) configEq(expected keyVal, fs afero.Fs) {
	a.t.Helper()
	f, e := fs.Open(TEST_CONF_FILE)
	a.ok(e)

	s, err := ioutil.ReadAll(f)
	a.ok(err)

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
