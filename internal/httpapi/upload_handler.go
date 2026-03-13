package httpapi

import (
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type UploadHandler struct {
	storage ports.StorageService
}

func NewUploadHandler(storage ports.StorageService) *UploadHandler {
	return &UploadHandler{storage: storage}
}

// @Summary Upload file
// @Tags admin
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/admin/upload [post]
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Validate content type (images only)
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		writeError(w, http.StatusBadRequest, "only image files are allowed")
		return
	}

	url, err := h.storage.Upload(r.Context(), header.Filename, contentType, header.Size, file)
	if err != nil {
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}
