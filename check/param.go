//Package check contains tools to validate the signatures for DDAPI requests
package check

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type (
	// ParamsChecker represents a param collection validator
	//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ParamsChecker
	ParamsChecker interface {
		Check(pm map[string]interface{}, schema string) error
	}

	// Params a params checker
	Params func(pm map[string]interface{}, s string) error
)

// Check checks if the params are valid
func (f Params) Check(pm map[string]interface{}, s string) error {
	return f(pm, s)
}

// JSONSchema a json-schema based params checker
func JSONSchema(pm map[string]interface{}, s string) error {
	if s == "" {
		return fmt.Errorf("params schema cannot be empty")
	}

	if len(pm) == 0 {
		return fmt.Errorf("params map cannot be nil or empty")
	}

	param := gojsonschema.NewGoLoader(pm)
	schema := gojsonschema.NewStringLoader(s)
	result, err := gojsonschema.Validate(schema, param)
	if err != nil {
		return fmt.Errorf("could not validate parameters: %w", err)
	}

	if !result.Valid() {
		return fmt.Errorf("invalid parameters %v", result.Errors())
	}

	return nil

}
