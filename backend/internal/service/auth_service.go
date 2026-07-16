package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo          ports.AuthRepository
	adminUser     string
	adminPassword string
	mode          string
	secret        []byte
	ttl           time.Duration
}

const (
	adminAuthModeDBThenEnv = "db_then_env"
	adminAuthModeDBOnly    = "db_only"
	adminAuthModeEnvOnly   = "env_only"
)

func NewAuthService(repo ports.AuthRepository, adminUser, adminPassword, jwtSecret, mode string, ttl time.Duration) (*AuthService, error) {
	adminUser = strings.TrimSpace(adminUser)
	jwtSecret = strings.TrimSpace(jwtSecret)
	mode = normalizeAdminAuthMode(mode)

	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if mode == adminAuthModeDBOnly && repo == nil {
		return nil, errors.New("ADMIN_AUTH_MODE=db_only requires auth repo")
	}
	if mode == adminAuthModeEnvOnly && (adminUser == "" || adminPassword == "") {
		return nil, errors.New("ADMIN_AUTH_MODE=env_only requires ADMIN_USER and ADMIN_PASSWORD")
	}
	if mode == adminAuthModeDBThenEnv && repo == nil && (adminUser == "" || adminPassword == "") {
		return nil, errors.New("admin creds are required when auth repo is not configured")
	}
	if ttl <= 0 {
		ttl = 8 * time.Hour
	}

	return &AuthService{
		repo:          repo,
		adminUser:     adminUser,
		adminPassword: adminPassword,
		mode:          mode,
		secret:        []byte(jwtSecret),
		ttl:           ttl,
	}, nil
}

var _ ports.AuthService = (*AuthService)(nil)

func (s *AuthService) Login(ctx context.Context, username, password string) (string, bool, error) {
	username = strings.TrimSpace(username)

	switch s.mode {
	case adminAuthModeDBOnly:
		return s.loginByDB(ctx, username, password)
	case adminAuthModeEnvOnly:
		return s.loginByEnv(username, password)
	default:
		if s.repo != nil {
			hash, ok, err := s.repo.GetAdminPasswordHash(ctx, username)
			if err != nil {
				return "", false, err
			}
			if ok {
				if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
					return "", false, nil
				}
				return s.signToken(username)
			}
		}
		return s.loginByEnv(username, password)
	}
}

func (s *AuthService) loginByDB(ctx context.Context, username, password string) (string, bool, error) {
	if s.repo == nil {
		return "", false, nil
	}
	hash, ok, err := s.repo.GetAdminPasswordHash(ctx, username)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", false, nil
	}
	return s.signToken(username)
}

func (s *AuthService) loginByEnv(username, password string) (string, bool, error) {
	if username != s.adminUser || password != s.adminPassword {
		return "", false, nil
	}
	return s.signToken(username)
}

func normalizeAdminAuthMode(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	switch v {
	case "", adminAuthModeDBThenEnv:
		return adminAuthModeDBThenEnv
	case adminAuthModeDBOnly:
		return adminAuthModeDBOnly
	case adminAuthModeEnvOnly:
		return adminAuthModeEnvOnly
	default:
		return adminAuthModeDBThenEnv
	}
}


func (s *AuthService) signToken(username string) (string, bool, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   username,
		Issuer:    "RadarTcell",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", false, err
	}
	return signed, true, nil
}

func (s *AuthService) Verify(ctx context.Context, token string) (string, bool, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false, nil
	}

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithIssuer("RadarTcell"))

	if err != nil || !parsed.Valid {
		return "", false, nil
	}
	if claims.Subject == "" {
		return "", false, nil
	}
	return claims.Subject, true, nil
}
