package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/at-silva/ddapi/session"
)

// ReadSession copies the session params from the session into the parameters collection
func ReadSession(s session.Reader, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, ok := r.Context().Value(DecodedParams).(map[string]interface{})
		if !ok {
			http.Error(w, errEncode(fmt.Errorf("could not copy session params: invalid params")), http.StatusInternalServerError)
			return
		}

		auth := r.Header.Get("Authorization")
		token := strings.Split(auth, "Bearer ")
		if len(token) != 2 {
			http.Error(w, errEncode(fmt.Errorf("could not copy session params: invalid Authorization header")), http.StatusBadRequest)
			return
		}

		err := s.Copy(token[1], params)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not copy session params: %w", err)), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
