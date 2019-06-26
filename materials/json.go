package materials

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
)

type Material interface {
	Type() string
	SetRequiredString(string, string) error
	SetString(string, string, string)
	SetStringIfFlagSet(string, string, *pflag.Flag)
	SetBoolIfFlagSet(string, bool, *pflag.Flag)
	SetBool(string, bool)

	SetAttributes(map[string]interface{}) error
}

func NewGit() *Git { return &Git{hash: hash{`auto_update`: true}} }
func NewHg() *Hg   { return &Hg{hash: hash{`auto_update`: true}} }
func NewSvn() *Svn { return &Svn{hash: hash{`auto_update`: true}} }
func NewP4() *P4   { return &P4{hash: hash{`auto_update`: true}} }
func NewTfs() *Tfs { return &Tfs{hash: hash{`auto_update`: true}} }

func FromMap(data map[string]interface{}) (Material, error) {
	r := hash(data)

	if _, ok := r[`type`]; !ok {
		return nil, fmt.Errorf(`Missing material "type" in JSON payload: %s`, data)
	}

	if t, ok := (r[`type`]).(string); !ok {
		return nil, fmt.Errorf(`"type" must be a string in material spec: %s`, data)
	} else {
		if _, ok := r[`attributes`]; !ok {
			return nil, fmt.Errorf(`Missing material "attributes" in JSON payload: %s`, data)
		}

		var attrs hash

		if v, ok := (r[`attributes`]).(map[string]interface{}); !ok {
			return nil, fmt.Errorf(`"attributes" must be a dict in material spec: %s`, data)
		} else {
			attrs = hash(v)
		}

		var mat Material
		switch t {
		case `git`:
			mat = NewGit()
		case `hg`:
			mat = NewHg()
		case `svn`:
			mat = NewSvn()
		case `p4`:
			mat = NewP4()
		case `tfs`:
			mat = NewTfs()
		default:
			return nil, fmt.Errorf(`Unknown material type %q`, t)
		}

		return mat, mat.SetAttributes(attrs)
	}
}

type Git struct {
	hash
}

func (g *Git) Type() string {
	return `git`
}

func (m *Git) SetAttributes(cfg map[string]interface{}) error {
	attrs := hash(cfg)
	if err := attrs.copyStrIfPresent(m.hash, `url`); err != nil {
		return err
	}
	if err := attrs.copyStrOrDefault(m.hash, `branch`, `master`); err != nil {
		return err
	}
	return attrs.copyStrIfNotNull(m.hash, `name`)
}

func (m *Git) asMap() hash {
	return hash{
		`type`:       m.Type(),
		`attributes`: m.hash.copyTo(make(hash)),
	}
}

func (m *Git) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.asMap())
}

func (m *Git) UnmarshalJSON(b []byte) error {
	if attrs, err := unmarshalAttrs(b, m.Type()); err != nil {
		return err
	} else {
		return m.SetAttributes(attrs)
	}
}

type Hg struct {
	hash
}

func (g *Hg) Type() string {
	return `hg`
}

func (m *Hg) SetAttributes(cfg map[string]interface{}) error {
	attrs := hash(cfg)
	if err := attrs.copyStrIfPresent(m.hash, `url`); err != nil {
		return err
	}
	return attrs.copyStrIfNotNull(m.hash, `name`)
}

func (m *Hg) asMap() hash {
	return hash{
		`type`:       m.Type(),
		`attributes`: m.hash.copyTo(make(hash)),
	}
}

func (m *Hg) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.asMap())
}

func (m *Hg) UnmarshalJSON(b []byte) error {
	if attrs, err := unmarshalAttrs(b, m.Type()); err != nil {
		return err
	} else {
		return m.SetAttributes(attrs)
	}
}

type Svn struct {
	hash
}

func (m *Svn) Type() string {
	return `svn`
}

func (m *Svn) SetAttributes(cfg map[string]interface{}) error {
	attrs := hash(cfg)
	if err := attrs.copyStrIfPresent(m.hash, `url`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `username`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `password`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `encrypted_password`); err != nil {
		return err
	}
	if err := attrs.copyBoolOrDefault(m.hash, `check_externals`, true); err != nil {
		return err
	}
	return attrs.copyStrIfNotNull(m.hash, `name`)
}

func (m *Svn) asMap() hash {
	return hash{
		`type`:       m.Type(),
		`attributes`: m.hash.copyTo(make(hash)),
	}
}

func (m *Svn) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.asMap())
}

func (m *Svn) UnmarshalJSON(b []byte) error {
	if attrs, err := unmarshalAttrs(b, m.Type()); err != nil {
		return err
	} else {
		return m.SetAttributes(attrs)
	}
}

type P4 struct {
	hash
}

func (m *P4) Type() string {
	return `p4`
}

func (m *P4) SetAttributes(cfg map[string]interface{}) error {
	attrs := hash(cfg)
	if err := attrs.copyStrIfPresent(m.hash, `port`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `view`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `username`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `password`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `encrypted_password`); err != nil {
		return err
	}
	if err := attrs.copyBoolOrDefault(m.hash, `use_tickets`, true); err != nil {
		return err
	}
	return attrs.copyStrIfNotNull(m.hash, `name`)
}

func (m *P4) asMap() hash {
	return hash{
		`type`:       m.Type(),
		`attributes`: m.hash.copyTo(make(hash)),
	}
}

func (m *P4) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.asMap())
}

func (m *P4) UnmarshalJSON(b []byte) error {
	if attrs, err := unmarshalAttrs(b, m.Type()); err != nil {
		return err
	} else {
		return m.SetAttributes(attrs)
	}
}

type Tfs struct {
	hash
}

func (m *Tfs) Type() string {
	return `tfs`
}

func (m *Tfs) SetAttributes(cfg map[string]interface{}) error {
	attrs := hash(cfg)
	if err := attrs.copyStrIfPresent(m.hash, `url`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `project_path`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `username`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `password`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `encrypted_password`); err != nil {
		return err
	}
	if err := attrs.copyStrIfPresent(m.hash, `domain`); err != nil {
		return err
	}
	return attrs.copyStrIfNotNull(m.hash, `name`)
}

func (m *Tfs) asMap() hash {
	return hash{
		`type`:       m.Type(),
		`attributes`: m.hash.copyTo(make(hash)),
	}
}

func (m *Tfs) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.asMap())
}

func (m *Tfs) UnmarshalJSON(b []byte) error {
	if attrs, err := unmarshalAttrs(b, m.Type()); err != nil {
		return err
	} else {
		return m.SetAttributes(attrs)
	}
}

func unmarshalAttrs(b []byte, materialType string) (hash, error) {
	attrs := make(hash)
	if err := json.Unmarshal(b, &attrs); err == nil {
		if t, err := attrs.string(`type`); err == nil {
			if t != materialType {
				return nil, fmt.Errorf(`expected material JSON "type" to be equal to %q but was actually %q`, materialType, t)
			}
		} else {
			return nil, err
		}

		return attrs.subHash(`attributes`)
	} else {
		return nil, err
	}
}
