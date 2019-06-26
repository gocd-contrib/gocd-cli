package configrepo_test

import (
	"encoding/json"
	"reflect"
	"testing"
)

type asserter struct {
	t *testing.T
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

func (a *asserter) jsonEq(expected, actual string) {
	a.t.Helper()

	j := make(map[string]interface{})
	k := make(map[string]interface{})

	if err := json.Unmarshal([]byte(expected), &j); err != nil {
		a.t.Errorf(`"expected" string is not valid JSON: %s`, expected)
	}

	if err := json.Unmarshal([]byte(actual), &k); err != nil {
		a.t.Errorf(`"actual" string is not valid JSON: %s`, actual)
	}

	if !reflect.DeepEqual(j, k) {
		a.t.Errorf("Expected:\n%s\n\nto be equivalent to:\n%s", expected, actual)
	}
}

func (a *asserter) deepEq(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Expected:\n%v\n\nto be deeply equal to:\n%v", expected, actual)
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
