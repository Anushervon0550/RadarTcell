package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/go-chi/chi/v5"
)

// ============================== JSON ENCODE =================================

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"error": msg})
}

// ============================== JSON DECODE =================================

const maxJSONBodyBytes int64 = 1 << 20 // 1 MiB

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return io.EOF
	}
	defer func() { _ = r.Body.Close() }()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("extra data after json")
	}

	return nil
}

func decodeJSONOr400(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	err := decodeJSON(r, dst)
	if err == nil {
		return true
	}
	var maxBytesErr *http.MaxBytesError

	switch {
	case errors.Is(err, io.EOF):
		writeError(w, http.StatusBadRequest, "empty body")
	case errors.As(err, &maxBytesErr):
		writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
	default:
		writeError(w, http.StatusBadRequest, "invalid json")
	}
	return false
}

// ============================== REQUEST PARAMS ==============================

func pathParamRequired(r *http.Request, name string) (string, bool) {
	v := strings.TrimSpace(chiURLParam(r, name))
	if v == "" {
		return "", false
	}
	return v, true
}

func queryParamRequired(r *http.Request, name string) (string, bool) {
	v := strings.TrimSpace(r.URL.Query().Get(name))
	if v == "" {
		return "", false
	}
	return v, true
}

func chiURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

// ============================== DOMAIN ERRORS ===============================

const internalErrorMessage = "internal server error"

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, internalErrorMessage)
}

func writeDomainErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeInternalError(w)
	}
}
