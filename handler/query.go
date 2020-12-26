package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/at-silva/ddapi/check"
	"github.com/at-silva/ddapi/db"
	"github.com/at-silva/ddapi/session"
)

type (
	queryHandler struct {
		db db.DB
	}

	queryResponse struct {
		Data  interface{} `json:"data"`
		Error *string     `json:"error"`
	}
)

// NewQuery returns a new DDApi query handler
func NewQuery(db db.DB, sc check.SignatureChecker, s session.Reader, pc check.ParamsChecker) http.Handler {
	h := DecodeRequest(
		CheckSignatures(sc,
			ReadSession(s,
				CheckParams(pc,
					queryHandler{
						db,
					}))))

	return h
}

func (h queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q, ok := r.Context().Value(DecodedRequest).(request)
	if !ok {
		http.Error(w, errEncode(fmt.Errorf("could not query the database: invalid request")), http.StatusInternalServerError)
		return
	}

	params, ok := r.Context().Value(DecodedParams).(map[string]interface{})
	if !ok {
		http.Error(w, errEncode(fmt.Errorf("could not query the database: invalid params")), http.StatusInternalServerError)
		return
	}

	rows, err := h.db.NamedQueryContext(r.Context(), q.SQL, params)
	if err != nil {
		http.Error(w, errEncode(fmt.Errorf("could not query the database: %w", err)), http.StatusInternalServerError)
		return
	}

	var res []interface{}
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)
		if err != nil {
			http.Error(w, errEncode(fmt.Errorf("could not scan rows: %w", err)), http.StatusInternalServerError)
			return
		}
		mapBytesToString(row)
		res = append(res, row)
	}

	resp, err := json.Marshal(queryResponse{res, nil})
	if err != nil {
		http.Error(w, errEncode(fmt.Errorf("could not serialize result: %w", err)), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, errEncode(fmt.Errorf("could not write to the response: %w", err)), http.StatusInternalServerError)
		return
	}
}
