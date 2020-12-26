/*Package handler contains a set of http handlers to address:
decode: DDAPI requests decoding
exec: DML execution
params: query/statement parameters validation
query: DQL execution
session: JWT/session introspection
signature: query/statement signature checking
*/
package handler

import (
	"encoding/json"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.Handler
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 io.Reader
type (
	request struct {
		SQL                   string `json:"sql"`
		SQLSignature          string `json:"sqlSignature"`
		Params                string `json:"params"`
		ParamsSchema          string `json:"paramsSchema"`
		ParamsSchemaSignature string `json:"paramsSchemaSignature"`
	}
	alias request
)

// UnmarshallJSON custom unmarshaller to allow lazy unmarshalling of ParamsSchema and Params fields
func (r *request) UnmarshalJSON(data []byte) error {
	aux := &struct {
		*alias
		Params       json.RawMessage `json:"params"`
		ParamsSchema json.RawMessage `json:"paramsSchema"`
	}{alias: (*alias)(r)}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	r.SQL = aux.SQL
	r.SQLSignature = aux.SQLSignature
	r.Params = string(aux.Params)
	r.ParamsSchema = string(aux.ParamsSchema)
	r.ParamsSchemaSignature = aux.ParamsSchemaSignature
	return nil
}

func mapBytesToString(s map[string]interface{}) {
	for k, v := range s {
		if b, ok := v.([]byte); ok {
			s[k] = string(b)
		}
	}
}

func errEncode(e error) string {
	res, _ := json.Marshal(&struct {
		Error string `json:"error"`
	}{
		Error: e.Error(),
	})

	return string(res)
}
