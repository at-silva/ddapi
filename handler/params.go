package handler

import (
	"fmt"
	"net/http"

	"github.com/at-silva/ddapi/check"
)

// CheckParams validates the parameters in an incoming request
func CheckParams(pc check.ParamsChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := r.Context().Value(DecodedParams).(map[string]interface{})
		if !ok {
			http.Error(w, errEncode(fmt.Errorf("could not check params: invalid params")), http.StatusInternalServerError)
			return
		}

		req, ok := r.Context().Value(DecodedRequest).(request)
		if !ok {
			http.Error(w, errEncode(fmt.Errorf("could not check params: invalid request")), http.StatusInternalServerError)
			return
		}

		err := pc.Check(p, req.ParamsSchema)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("invalid params: %w", err)), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
