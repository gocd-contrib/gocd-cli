package api

import (
	"fmt"
	"strings"
)

type CrError struct {
	File string `json:"location"`
	Msg  string `json:"message"`
}

func (cr *CrError) HumanFormat() string {
	return fmt.Sprintf("%s:\n  %s", cr.File, cr.Msg)
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
