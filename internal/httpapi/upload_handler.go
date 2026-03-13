package httpapi

import (
	"io"
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
	defer func() { _ = file.Close() }()

	// Validate real content type from file signature (first 512 bytes).
	buf := make([]byte, 512)
	n, readErr := io.ReadFull(file, buf)
	if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
		writeError(w, http.StatusBadRequest, "failed to read uploaded file")
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeError(w, http.StatusBadRequest, "failed to process uploaded file")
		return
	}

	contentType := http.DetectContentType(buf[:n])
	if !strings.HasPrefix(contentType, "image/") {
		writeError(w, http.StatusBadRequest, "only image files are allowed")
		return
	}

	url, err := h.storage.Upload(r.Context(), header.Filename, contentType, -1, file)
	if err != nil {
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}
