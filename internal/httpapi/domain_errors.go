package httpapi

import (
	"errors"
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

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
