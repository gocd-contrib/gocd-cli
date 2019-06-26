package configrepo

import (
	"encoding/json"
	"testing"

	"github.com/gocd-contrib/gocd-cli/materials"
)

const (
	data = `{
    "id": "repo-2",
    "plugin_id": "json.config.plugin",
    "material": {
      "type": "hg",
      "attributes": {
        "url": "https://hgbucket.org/repo",
        "auto_update": true
      }
    },
    "configuration": [
      {
        "key": "pattern",
        "value": "*.myextension"
      },
      {
        "key": "token",
        "encrypted_value": "abcd1234"
      }
    ]
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

func BenchmarkUnmarshalConfigRepo(b *testing.B) {
	j := []byte(data)

	for i := 0; i < b.N; i++ {
		v := &ConfigRepo{}
		if err := json.Unmarshal(j, &v); err != nil {
			b.Errorf(`Failed to deserialize JSON: %v`, err)
		}
	}
}

func BenchmarkMarshalConfigRepo(b *testing.B) {
	repo := NewConfigRepo(
		`repo-2`,
		`json.config.plugin`,
		materials.NewHg(),
		NewPlainTextProperty(`pattern`, `*.myextension`),
		NewSecretProperty(`token`, `abcd1234`),
	)

	repo.Material.SetAttributes(map[string]interface{}{
		`url`: `https://hgbucket.org/repo`,
	})

	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(repo); err != nil {
			b.Errorf(`Failed to serialize ConfigRepo: %v`, err)
		}
	}
}

func BenchmarkEquiv(b *testing.B) {
	var x Property = &PlainTextProperty{Key: "foo", Value: "bar"}
	var y Property = &PlainTextProperty{Key: "foo", Value: "bar"}

	for i := 0; i < b.N; i++ {
		x.Equivalent(y)
	}
}

var drain bool

func BenchmarkEqual(b *testing.B) {
	var x = PlainTextProperty{Key: "foo", Value: "bar"}
	var y = PlainTextProperty{Key: "foo", Value: "bar"}

	for i := 0; i < b.N; i++ {
		drain = x == y
	}
}

func BenchmarkEqualDeref(b *testing.B) {
	var x = &PlainTextProperty{Key: "foo", Value: "bar"}
	var y = &PlainTextProperty{Key: "foo", Value: "bar"}

	for i := 0; i < b.N; i++ {
		drain = *x == *y
	}
}
