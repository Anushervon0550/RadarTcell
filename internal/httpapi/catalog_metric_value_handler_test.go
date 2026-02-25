package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
)

type catalogMetricValueServiceStub struct {
	ports.CatalogService // embedding: не нужно реализовывать весь интерфейс

	getMetricValueFn func(ctx context.Context, metricID, techID string) (map[string]any, bool, error)

	gotMetricID string
	gotTechID   string
}

func (s *catalogMetricValueServiceStub) GetMetricValue(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
	s.gotMetricID = metricID
	s.gotTechID = techID

	if s.getMetricValueFn != nil {
		return s.getMetricValueFn(ctx, metricID, techID)
	}

	return nil, false, nil
}

func withChiMetricID(r *http.Request, metricID string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", metricID)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestCatalogHandler_GetMetricValue_OK(t *testing.T) {
	stub := &catalogMetricValueServiceStub{
		getMetricValueFn: func(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
			return map[string]any{
				"metric_id":     metricID,
				"technology_id": techID,
				"metric_name":   "Custom Metric 01",
				"field_key":     "custom_metric_1",
				"type":          "bubble",
				"value":         0.7,
			}, true, nil
		},
	}

	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m1/values?technology_id=t1", nil)
	req = withChiMetricID(req, "m1")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if stub.gotMetricID != "m1" {
		t.Fatalf("expected metricID=m1, got %q", stub.gotMetricID)
	}
	if stub.gotTechID != "t1" {
		t.Fatalf("expected techID=t1, got %q", stub.gotTechID)
	}

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v; body=%s", err, rr.Body.String())
	}

	if got["metric_id"] != "m1" {
		t.Fatalf("expected metric_id=m1, got %v", got["metric_id"])
	}
	if got["technology_id"] != "t1" {
		t.Fatalf("expected technology_id=t1, got %v", got["technology_id"])
	}
	if got["field_key"] != "custom_metric_1" {
		t.Fatalf("expected field_key=custom_metric_1, got %v", got["field_key"])
	}
}

func TestCatalogHandler_GetMetricValue_BadRequest_MissingTechnologyID(t *testing.T) {
	stub := &catalogMetricValueServiceStub{}
	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m1/values", nil)
	req = withChiMetricID(req, "m1")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}

	if !strings.Contains(strings.ToLower(rr.Body.String()), "technology_id is required") {
		t.Fatalf("expected 'technology_id is required', got: %s", rr.Body.String())
	}
}

func TestCatalogHandler_GetMetricValue_BadRequest_MissingMetricID(t *testing.T) {
	stub := &catalogMetricValueServiceStub{}
	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics//values?technology_id=t1", nil)
	// id в chi route params специально НЕ добавляем
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}

	if !strings.Contains(strings.ToLower(rr.Body.String()), "metric id is required") {
		t.Fatalf("expected 'metric id is required', got: %s", rr.Body.String())
	}
}

func TestCatalogHandler_GetMetricValue_NotFound(t *testing.T) {
	stub := &catalogMetricValueServiceStub{
		getMetricValueFn: func(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
			return nil, false, nil
		},
	}

	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m404/values?technology_id=t1", nil)
	req = withChiMetricID(req, "m404")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
	}

	if !strings.Contains(strings.ToLower(rr.Body.String()), "not found") {
		t.Fatalf("expected 'not found', got: %s", rr.Body.String())
	}
}

func TestCatalogHandler_GetMetricValue_DomainInvalid_To400(t *testing.T) {
	stub := &catalogMetricValueServiceStub{
		getMetricValueFn: func(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
			return nil, false, fmt.Errorf("%w: bad input", domain.ErrInvalid)
		},
	}

	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m1/values?technology_id=t1", nil)
	req = withChiMetricID(req, "m1")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCatalogHandler_GetMetricValue_DomainConflict_To409(t *testing.T) {
	stub := &catalogMetricValueServiceStub{
		getMetricValueFn: func(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
			return nil, false, fmt.Errorf("%w: conflict", domain.ErrConflict)
		},
	}

	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m1/values?technology_id=t1", nil)
	req = withChiMetricID(req, "m1")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusConflict, rr.Code, rr.Body.String())
	}
}

func TestCatalogHandler_GetMetricValue_UnknownError_To500(t *testing.T) {
	stub := &catalogMetricValueServiceStub{
		getMetricValueFn: func(ctx context.Context, metricID, techID string) (map[string]any, bool, error) {
			return nil, false, fmt.Errorf("db exploded")
		},
	}

	h := NewCatalogHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/m1/values?technology_id=t1", nil)
	req = withChiMetricID(req, "m1")
	rr := httptest.NewRecorder()

	h.GetMetricValue(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}
