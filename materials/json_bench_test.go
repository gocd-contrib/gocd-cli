package materials

import (
	"encoding/json"
	"testing"
)

const (
	data = `{
    "type": "tfs",
    "attributes": {
      "url": "http://foobar.com",
      "project_path": "my-project_path",
      "domain": "my-domain",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz",
      "auto_update": true
    }
  }`
)

func BenchmarkUnmarshalGenericMapAsBaseline(b *testing.B) {
	j := []byte(data)

	for i := 0; i < b.N; i++ {
		v := make(map[string]interface{})
		if err := json.Unmarshal(j, &v); err != nil {
			b.Errorf(`Failed to deserialize JSON: %v`, err)
		}
	}
}

// This is about twice as slow as the UnmarshalJSON() method.
// encoding/json must be doing something expensive before hitting
// UnmarshalJSON() on the Material.
func BenchmarkUnmarshalMaterial(b *testing.B) {
	j := []byte(data)
	for i := 0; i < b.N; i++ {
		tfs := NewTfs()
		if err := json.Unmarshal(j, &tfs); err != nil {
			b.Errorf(`Cannot parse JSON: %v`, err)
		}
	}
}

// Bench the Unmarshaler interface implementation
func Benchmark_Material_UnmarshalJSON(b *testing.B) {
	j := []byte(data)
	for i := 0; i < b.N; i++ {
		tfs := NewTfs()
		if err := tfs.UnmarshalJSON(j); err != nil {
			b.Errorf(`Cannot parse JSON: %v`, err)
		}
	}
}

// This is about 50% as slower than the MarshalJSON() method.
// encoding/json must be doing something expensive before hitting
// MarshalJSON() on the Material.
func BenchmarkMarshalMaterial(b *testing.B) {
	m := NewTfs()
	if err := m.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	}); err != nil {
		b.Errorf(`Cannot build Material: %v`, err)
	}

	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(m); err != nil {
			b.Errorf(`Cannot serialize Material: %v`, err)
		}
	}
}

// Bench the Marshaler interface implementation
func Benchmark_Material_MarshalJSON(b *testing.B) {
	m := NewTfs()
	if err := m.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	}); err != nil {
		b.Errorf(`Cannot build Material: %v`, err)
	}

	for i := 0; i < b.N; i++ {
		if _, err := m.MarshalJSON(); err != nil {
			b.Errorf(`Cannot serialize Material: %v`, err)
		}
	}
}

func BenchmarkEquivalent(b *testing.B) {
	this := NewTfs()
	if err := this.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	}); err != nil {
		b.Errorf(`Cannot build Material: %v`, err)
	}

	that := NewTfs()
	if err := that.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	}); err != nil {
		b.Errorf(`Cannot build Material: %v`, err)
	}

	for i := 0; i < b.N; i++ {
		this.Equivalent(that)
	}
}

func Benchmark_unmarshalAttrs(b *testing.B) {
	j := []byte(data)
	for i := 0; i < b.N; i++ {
		if _, err := unmarshalAttrs(j, `tfs`); err != nil {
			b.Errorf(`Failed to deserialize JSON: %v`, err)
		}
	}
}

func Benchmark_SetAttributes(b *testing.B) {
	sampleData := map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	}
	m := NewTfs()
	for i := 0; i < b.N; i++ {
		if err := m.SetAttributes(sampleData); err != nil {
			b.Errorf(`Cannot build Material: %v`, err)
		}
	}
}
