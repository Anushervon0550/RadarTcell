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

type preferencesServiceStub struct {
	ports.PreferencesService
	saveFn func(ctx context.Context, p domain.Preferences) error
	getFn  func(ctx context.Context, userID string) (domain.Preferences, bool, error)

	gotSave domain.Preferences
	gotGet  string
}

func (s *preferencesServiceStub) Save(ctx context.Context, p domain.Preferences) error {
	s.gotSave = p
	if s.saveFn != nil {
		return s.saveFn(ctx, p)
	}
	return nil
}

func (s *preferencesServiceStub) Get(ctx context.Context, userID string) (domain.Preferences, bool, error) {
	s.gotGet = userID
	if s.getFn != nil {
		return s.getFn(ctx, userID)
	}
	var zero domain.Preferences
	return zero, false, nil
}

func withChiUserID(r *http.Request, userID string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", userID)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestPreferencesHandler_Save_OK(t *testing.T) {
	stub := &preferencesServiceStub{}
	h := NewPreferencesHandler(stub)

	body := `{"user_id":"u1","settings":{"show_tail":true,"magnetic_label":false}}`
	req := httptest.NewRequest(http.MethodPost, "/api/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Save(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if strings.TrimSpace(stub.gotSave.UserID) != "u1" {
		t.Fatalf("expected saved user_id=u1, got %q", stub.gotSave.UserID)
	}

	if len(stub.gotSave.Settings) == 0 {
		t.Fatalf("expected non-empty settings json")
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v, body=%s", err, rr.Body.String())
	}
	if resp["status"] != "ok" {
		t.Fatalf("expected response status=ok, got %v", resp["status"])
	}
}

func TestPreferencesHandler_Save_BadJSON(t *testing.T) {
	stub := &preferencesServiceStub{}
	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodPost, "/api/preferences", strings.NewReader(`{"user_id":`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestPreferencesHandler_Save_DomainInvalid_To400(t *testing.T) {
	stub := &preferencesServiceStub{
		saveFn: func(ctx context.Context, p domain.Preferences) error {
			return fmt.Errorf("%w: user id is required", domain.ErrInvalid)
		},
	}
	h := NewPreferencesHandler(stub)

	body := `{"user_id":"","settings":{"show_tail":true}}`
	req := httptest.NewRequest(http.MethodPost, "/api/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Save(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestPreferencesHandler_Get_OK(t *testing.T) {
	stub := &preferencesServiceStub{
		getFn: func(ctx context.Context, userID string) (domain.Preferences, bool, error) {
			return domain.Preferences{
				UserID:   userID,
				Settings: json.RawMessage(`{"show_tail":true,"magnetic_label":false}`),
			}, true, nil
		},
	}

	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/u1", nil)
	req = withChiUserID(req, "u1")
	rr := httptest.NewRecorder()

	h.Get(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if stub.gotGet != "u1" {
		t.Fatalf("expected userID=u1 passed to service, got %q", stub.gotGet)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v, body=%s", err, rr.Body.String())
	}

	if resp["user_id"] != "u1" {
		t.Fatalf("expected response user_id=u1, got %v", resp["user_id"])
	}

	settings, ok := resp["settings"].(map[string]any)
	if !ok {
		t.Fatalf("expected settings object, got %T (%v)", resp["settings"], resp["settings"])
	}

	if _, exists := settings["show_tail"]; !exists {
		t.Fatalf("expected settings.show_tail in response")
	}
}

func TestPreferencesHandler_Get_BadRequest_MissingUserID(t *testing.T) {
	stub := &preferencesServiceStub{}
	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/", nil)
	// user_id в chi params специально НЕ добавляем
	rr := httptest.NewRecorder()

	h.Get(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}

	if !strings.Contains(strings.ToLower(rr.Body.String()), "user id is required") {
		t.Fatalf("expected 'user id is required', got: %s", rr.Body.String())
	}
}

func TestPreferencesHandler_Get_NotFound(t *testing.T) {
	stub := &preferencesServiceStub{
		getFn: func(ctx context.Context, userID string) (domain.Preferences, bool, error) {
			var zero domain.Preferences
			return zero, false, nil
		},
	}

	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/u404", nil)
	req = withChiUserID(req, "u404")
	rr := httptest.NewRecorder()

	h.Get(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
	}

	if !strings.Contains(strings.ToLower(rr.Body.String()), "not found") {
		t.Fatalf("expected 'not found', got: %s", rr.Body.String())
	}
}

func TestPreferencesHandler_Get_DomainInvalid_To400(t *testing.T) {
	stub := &preferencesServiceStub{
		getFn: func(ctx context.Context, userID string) (domain.Preferences, bool, error) {
			var zero domain.Preferences
			return zero, false, fmt.Errorf("%w: bad user id", domain.ErrInvalid)
		},
	}

	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/u1", nil)
	req = withChiUserID(req, "u1")
	rr := httptest.NewRecorder()

	h.Get(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}
