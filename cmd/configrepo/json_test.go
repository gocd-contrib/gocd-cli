package configrepo_test

import (
	"encoding/json"
	"testing"

	"github.com/gocd-contrib/gocd-cli/cmd/configrepo"
	"github.com/gocd-contrib/gocd-cli/materials"
)

func TestUnmarshalConfigRepo(t *testing.T) {
	as := asserts(t)
	actual := &configrepo.ConfigRepo{}

	as.ok(json.Unmarshal([]byte(`{
	  "id": "repo-2",
	  "plugin_id": "json.config.plugin",
	  "material": {
	    "type": "git",
	    "attributes": {
	      "url": "https://github.com/config-repo/gocd-json-config-example2.git",
	      "branch": "master",
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
	}`), &actual))

	expected := configrepo.NewConfigRepo(
		`repo-2`,
		`json.config.plugin`,
		materials.NewGit(),
		configrepo.NewPlainTextProperty(`pattern`, `*.myextension`),
		configrepo.NewSecretProperty(`token`, `abcd1234`),
	)

	expected.Material.SetAttributes(map[string]interface{}{
		`url`:    `https://github.com/config-repo/gocd-json-config-example2.git`,
		`branch`: `master`,
	})

	as.deepEq(expected, actual)
}

func TestMarshalConfigRepo(t *testing.T) {
	as := asserts(t)

	repo := configrepo.NewConfigRepo(
		`repo-2`,
		`json.config.plugin`,
		materials.NewGit(),
		configrepo.NewPlainTextProperty(`pattern`, `*.myextension`),
		configrepo.NewSecretProperty(`token`, `abcd1234`),
	)

	repo.Material.SetAttributes(map[string]interface{}{
		`url`:    `https://github.com/config-repo/gocd-json-config-example2.git`,
		`branch`: `master`,
	})

	b, err := json.MarshalIndent(repo, ``, `  `)
	as.ok(err)

	as.jsonEq(`{
	  "id": "repo-2",
	  "plugin_id": "json.config.plugin",
	  "material": {
	    "type": "git",
	    "attributes": {
	      "url": "https://github.com/config-repo/gocd-json-config-example2.git",
	      "branch": "master",
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
	}`, string(b))
}
