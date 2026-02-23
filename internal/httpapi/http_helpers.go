package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Запрещаем второй JSON-объект в body
	if err := dec.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
