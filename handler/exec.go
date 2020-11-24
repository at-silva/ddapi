package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/at-silva/ddapi/session"

	"github.com/at-silva/ddapi/check"
	"github.com/at-silva/ddapi/db"
)

type (
	execHandler struct {
		db db.DB
	}

	execResponse struct {
		RowsAffected   int64   `json:"rowsAffected"`
		LastInsertedID int64   `json:"lastInsertedId"`
		Error          *string `json:"error"`
	}
)

// NewExec returns a new DDApi exec handler
func NewExec(db db.DB, sc check.SignatureChecker, s session.Reader, pc check.ParamsChecker) http.Handler {
	h := DecodeRequest(
		CheckSignatures(sc,
			ReadSession(s,
				CheckParams(pc,
					execHandler{
						db,
					}))))

	return h
}

func execError(err error) string {
	e := err.Error()
	resp, _ := json.Marshal(execResponse{Error: &e})
	return string(resp)
}

func (h execHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q, ok := r.Context().Value(DecodedRequest).(request)
	if !ok {
		http.Error(w, execError(fmt.Errorf("could not query the database: invalid request")), http.StatusInternalServerError)
		return
	}

	params, ok := r.Context().Value(DecodedParams).(map[string]interface{})
	if !ok {
		http.Error(w, execError(fmt.Errorf("could not query the database: invalid params")), http.StatusInternalServerError)
		return
	}

	res, err := h.db.NamedExecContext(r.Context(), q.SQL, params)
	if err != nil {
		http.Error(w, execError(fmt.Errorf("could not query the database: %w", err)), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, execError(fmt.Errorf("could not read the number of rows affected: %w", err)), http.StatusInternalServerError)
		return
	}

	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, execError(fmt.Errorf("could not read the last inserted id: %w", err)), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(execResponse{rowsAffected, lastInsertedID, nil})
	if err != nil {
		log.Printf("marshal failed: %v", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
