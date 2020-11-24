package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type contextKey int

// DecodedRequest context key
const (
	DecodedRequest contextKey = iota
	DecodedParams
)

// DecodeRequest decodes an incoming request and adds it to the context
func DecodeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var q request

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not read body: %w", err)), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(b, &q)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not unmarshal body: %w", err)), http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), DecodedRequest, q))

		p := map[string]interface{}{}
		err = json.Unmarshal([]byte(q.Params), &p)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not unmarshal params: %w", err)), http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), DecodedParams, p))

		next.ServeHTTP(w, r)
	})
}
