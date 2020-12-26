package handler

import (
	"fmt"
	"net/http"
)

// ValidateFormRequest validates an incoming request
func ValidateFormRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		//TODO: implement character allow list
		if r.FormValue("sql") == "" {
			http.Error(w, errEncode(fmt.Errorf("invalid request: sql is required")), http.StatusBadRequest)
			return
		}

		if r.FormValue("sqlSignature") == "" {
			http.Error(w, errEncode(fmt.Errorf("invalid request: sqlSignature is required")), http.StatusBadRequest)
			return
		}

		if r.FormValue("params") == "" {
			http.Error(w, errEncode(fmt.Errorf("invalid request: params is required")), http.StatusBadRequest)
			return
		}

		if r.FormValue("paramsSchema") == "" {
			http.Error(w, errEncode(fmt.Errorf("invalid request: paramsSchema is required")), http.StatusBadRequest)
			return
		}

		if r.FormValue("paramsSchemaSignature") == "" {
			http.Error(w, errEncode(fmt.Errorf("invalid request: paramsSchemaSignature is required")), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
