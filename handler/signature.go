package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/at-silva/ddapi/check"
)

// CheckSignatures checks the signatures for a given request
func CheckSignatures(sc check.SignatureChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(DecodedRequest).(request)
		if !ok {
			http.Error(w, errEncode(fmt.Errorf("could not check signatures: invalid request")), http.StatusInternalServerError)
			return
		}

		s, err := base64.StdEncoding.DecodeString(req.SQLSignature)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not decode sql signature: %w", err)), http.StatusBadRequest)
			return
		}

		err = sc.Check([]byte(req.SQL), s)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not validate sql signature: %w", err)), http.StatusForbidden)
			return
		}

		s, err = base64.StdEncoding.DecodeString(req.ParamsSchemaSignature)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not decode params schema signature: %w", err)), http.StatusBadRequest)
			return
		}

		err = sc.Check([]byte(req.ParamsSchema), s)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not validate params schema signature: %w", err)), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
