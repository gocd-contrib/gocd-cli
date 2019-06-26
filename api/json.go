package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
)

// Creates an `io.ReadCloser` (i.e., `io.Reader` + `io.Closer`) from
// any serializable `interface{}` that returns the marshaled JSON output
// when calling `Read([]byte) (int, error)`.
func JsonReader(thing interface{}) io.ReadCloser {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := json.NewEncoder(w).Encode(thing); err != nil {
			w.CloseWithError(err)
		}
	}()
	return r
}

func PrettyPrintJson(thing interface{}) error {
	if b, err := json.MarshalIndent(thing, ``, `  `); err == nil {
		utils.Echofln(string(b))
		return nil
	} else {
		return err
	}
}

type MessageResponse interface {
	String() string
}

type CrError struct {
	File string `json:"location"`
	Msg  string `json:"message"`
}

func (cr *CrError) HumanFormat() string {
	return fmt.Sprintf("%s:\n  - %s", cr.File, cr.Msg)
}

type CrResponse struct {
	Errors []CrError `json:"errors"`
}

func (cr *CrResponse) DisplayErrors() string {
	strs := []string{}
	for _, es := range cr.Errors {
		strs = append(strs, es.HumanFormat())
	}
	return strings.Join(strs, "\n\n")
}

type CrPreflightResponse struct {
	Errors []string `json:"errors"`
	Valid  bool     `json:"valid"`
}

func (cr *CrPreflightResponse) DisplayErrors() string {
	return strings.Join(cr.Errors, "\n\n")
}

type BasicMessage struct {
	Message string `json:"message"`
}

func (m *BasicMessage) String() string {
	return m.Message
}

func ParseMessage(body []byte) (MessageResponse, error) {
	m := &BasicMessage{}

	if err := json.Unmarshal(body, m); err == nil {
		return m, nil
	} else {
		return nil, err
	}
}

func ParseCrMessageWithErrors(body []byte) (MessageResponse, error) {
	m := &CrMessageWithErrors{}

	if err := json.Unmarshal(body, m); err == nil {
		return m, nil
	} else {
		return nil, err
	}
}

type CrMessageWithErrors struct {
	Message string `json:"message"`
	Data    *struct {
		Material *struct {
			Attributes *struct {
				Errors ErrorMap `json:"errors"`
			} `json:"attributes"`
		} `json:"material"`
		Errors ErrorMap `json:"errors"`
	} `json:"data"`
}

func (c *CrMessageWithErrors) String() string {
	s := &strings.Builder{}
	s.WriteString(c.Message)
	if c.Data != nil {
		if len(c.Data.Errors) > 0 {
			s.WriteString(c.Data.Errors.String())
		}
		if c.Data.Material != nil && c.Data.Material.Attributes != nil {
			if len(c.Data.Material.Attributes.Errors) > 0 {
				s.WriteString(c.Data.Material.Attributes.Errors.String())
			}
		}
	} else {
		s.WriteString("\n")
	}
	return s.String()
}

type ErrorMap map[string][]string

func (m ErrorMap) String() string {
	s := &strings.Builder{}
	for k, v := range m {
		s.WriteString("\n\n")
		s.WriteString(`  ` + k + `: `)
		spaces := len(k) + 4

		for i, msg := range v {
			if i == 0 {
				s.WriteString(msg)
			} else {
				s.WriteString("\n" + strings.Repeat(` `, spaces) + msg)
			}
		}
	}
	return s.String()
}
