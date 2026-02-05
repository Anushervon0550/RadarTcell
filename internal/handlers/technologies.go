package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TechnologyDTO struct {
	ID             string `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Index          int    `json:"index"`
	ReadinessLevel int    `json:"readiness_level"`
	TrendSlug      string `json:"trend_slug"`
}

func (h *Handlers) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Фильтр: trend=ai
	trend := r.URL.Query().Get("trend")
	var trendArg any
	if trend == "" {
		trendArg = nil
	} else {
		trendArg = trend
	}

	rows, err := h.pool.Query(ctx, `
		SELECT tech.id, tech.slug, tech.name, tech."index", tech.readiness_level, t.slug
		FROM technologies tech
		JOIN trends t ON t.id = tech.trend_id
		WHERE ($1::text IS NULL OR t.slug = $1)
		ORDER BY t.order_index, tech."index"
	`, trendArg)
	if err != nil {
		http.Error(w, "db query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []TechnologyDTO
	for rows.Next() {
		var x TechnologyDTO
		if err := rows.Scan(&x.ID, &x.Slug, &x.Name, &x.Index, &x.ReadinessLevel, &x.TrendSlug); err != nil {
			http.Error(w, "db scan error", http.StatusInternalServerError)
			return
		}
		out = append(out, x)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}
