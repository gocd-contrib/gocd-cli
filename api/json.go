package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

type ApiMessage struct {
	Message string `json:"message"`
}

func (am *ApiMessage) String() string {
	return am.Message
}

func ParseMessage(body []byte) (*ApiMessage, error) {
	m := &ApiMessage{}

	if err := json.Unmarshal(body, m); err == nil {
		return m, nil
	} else {
		return nil, err
	}
}
