package materials_test

import (
	"encoding/json"
	"testing"

	"github.com/gocd-contrib/gocd-cli/materials"
)

func TestGitMarshal(t *testing.T) {
	as := asserts(t)
	m := materials.NewGit()
	m.SetAttributes(map[string]interface{}{
		`url`:    `http://foobar.com`,
		`branch`: `twig`,
	})

	b, err := json.MarshalIndent(m, ``, `  `)
	as.ok(err)

	as.jsonEq(string(b), `{
    "type": "git",
    "attributes": {
      "url": "http://foobar.com",
      "branch": "twig",
      "auto_update": true
    }
  }`)
}

func TestGitUnmarshal(t *testing.T) {
	as := asserts(t)
	expected := materials.NewGit()
	expected.SetAttributes(map[string]interface{}{
		`url`:         `http://foobar.com`,
		`branch`:      `twig`,
		`auto_update`: true,
	})

	actual := materials.NewGit()
	as.ok(json.Unmarshal([]byte(`{
    "type": "git",
    "attributes": {
      "url": "http://foobar.com",
      "branch": "twig",
      "auto_update": true
    }
  }`), &actual))

	as.deepEq(expected, actual)
}

func TestHgMarshal(t *testing.T) {
	as := asserts(t)
	m := materials.NewHg()
	m.SetAttributes(map[string]interface{}{
		`url`: `http://foobar.com`,
	})

	b, err := json.MarshalIndent(m, ``, `  `)
	as.ok(err)

	as.jsonEq(`{
    "type": "hg",
    "attributes": {
      "url": "http://foobar.com",
      "auto_update": true
    }
  }`, string(b))
}

func TestHgUnmarshal(t *testing.T) {
	as := asserts(t)
	expected := materials.NewHg()
	expected.SetAttributes(map[string]interface{}{
		`url`:         `http://foobar.com`,
		`branch`:      `twig`,
		`auto_update`: true,
	})

	actual := materials.NewHg()
	as.ok(json.Unmarshal([]byte(`{
    "type": "hg",
    "attributes": {
      "url": "http://foobar.com",
      "auto_update": true
    }
  }`), &actual))

	as.deepEq(expected, actual)
}

func TestSvnMarshal(t *testing.T) {
	as := asserts(t)
	m := materials.NewSvn()
	m.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
		`check_externals`:    true,
	})

	b, err := json.MarshalIndent(m, ``, `  `)
	as.ok(err)

	as.jsonEq(string(b), `{
    "type": "svn",
    "attributes": {
      "url": "http://foobar.com",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz",
      "check_externals": true,
      "auto_update": true
    }
  }`)
}

func TestSvnUnmarshal(t *testing.T) {
	as := asserts(t)
	expected := materials.NewSvn()
	expected.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
		`check_externals`:    true,
		`auto_update`:        true,
	})

	actual := materials.NewSvn()
	as.ok(json.Unmarshal([]byte(`{
    "type": "svn",
    "attributes": {
      "url": "http://foobar.com",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz",
      "check_externals": true
    }
  }`), &actual))

	as.deepEq(expected, actual)
}

func TestP4Marshal(t *testing.T) {
	as := asserts(t)
	m := materials.NewP4()
	m.SetAttributes(map[string]interface{}{
		`port`:               `foobar.com:443`,
		`view`:               `my-view`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
		`use_tickets`:        true,
	})

	b, err := json.MarshalIndent(m, ``, `  `)
	as.ok(err)

	as.jsonEq(string(b), `{
    "type": "p4",
    "attributes": {
      "port": "foobar.com:443",
      "view": "my-view",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz",
      "use_tickets": true,
      "auto_update": true
    }
  }`)
}

func TestP4Unmarshal(t *testing.T) {
	as := asserts(t)
	expected := materials.NewP4()
	expected.SetAttributes(map[string]interface{}{
		`port`:               `foobar.com:443`,
		`view`:               `my-view`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
		`use_tickets`:        true,
		`auto_update`:        true,
	})

	actual := materials.NewP4()
	as.ok(json.Unmarshal([]byte(`{
    "type": "p4",
    "attributes": {
      "port": "foobar.com:443",
      "view": "my-view",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz",
      "use_tickets": true
    }
  }`), &actual))

	as.deepEq(expected, actual)
}

func TestTfsMarshal(t *testing.T) {
	as := asserts(t)
	m := materials.NewTfs()
	m.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
	})

	b, err := json.MarshalIndent(m, ``, `  `)
	as.ok(err)

	as.jsonEq(string(b), `{
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
  }`)
}

func TestTfsUnmarshal(t *testing.T) {
	as := asserts(t)
	expected := materials.NewTfs()
	expected.SetAttributes(map[string]interface{}{
		`url`:                `http://foobar.com`,
		`project_path`:       `my-project_path`,
		`domain`:             `my-domain`,
		`username`:           `admin`,
		`password`:           `foo`,
		`encrypted_password`: `baz`,
		`auto_update`:        true,
	})

	actual := materials.NewTfs()
	as.ok(json.Unmarshal([]byte(`{
    "type": "tfs",
    "attributes": {
      "url": "http://foobar.com",
      "project_path": "my-project_path",
      "domain": "my-domain",
      "username": "admin",
      "password": "foo",
      "encrypted_password": "baz"
    }
  }`), &actual))

	as.deepEq(expected, actual)
}
