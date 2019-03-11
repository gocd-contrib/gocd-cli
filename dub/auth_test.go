package dub

import (
	"testing"
)

func TestBasicAuth(t *testing.T) {
	a := NewBasicAuth("foo", "bar")

	if b, ok := a.(*BasicAuth); ok {
		if "foo:bar" != string(b.payload()) {
			t.Errorf("BasicAuth should concat user and pass with a colon")
		}
	} else {
		t.Errorf("Expected a *BasicAuth, but got %T instead", b)
	}

	if "Basic Zm9vOmJhcg==" != a.Token() {
		t.Errorf("BasicAuth should output auth type with base64 payload")
	}
}

func TestTokenAuth(t *testing.T) {
	a := NewTokenAuth("abc123")

	if "Bearer abc123" != a.Token() {
		t.Errorf("TokenAuth should output auth type token payload")
	}
}
