package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TrendDTO struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	OrderIndex  int    `json:"order_index"`
}

func (h *Handlers) GetTrends(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx, `
		SELECT id, slug, name, COALESCE(description,''), COALESCE(image_url,''), order_index
		FROM trends
		ORDER BY order_index, slug
	`)
	if err != nil {
		http.Error(w, "db query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []TrendDTO
	for rows.Next() {
		var t TrendDTO
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.ImageURL, &t.OrderIndex); err != nil {
			http.Error(w, "db scan error", http.StatusInternalServerError)
			return
		}
		out = append(out, t)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}
