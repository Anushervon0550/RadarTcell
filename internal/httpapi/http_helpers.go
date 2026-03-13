package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxJSONBodyBytes int64 = 1 << 20 // 1 MiB

// Совместимый helper (старый стиль): handlers могут делать
// if err := decodeJSON(r, &req); err != nil { ... }
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

	// Запрещаем мусор после JSON
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("extra data after json")
	}

	return nil
}

// Новый удобный helper (если хочешь использовать в новых handlers):
// if !decodeJSONOr400(w, r, &req) { return }
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
