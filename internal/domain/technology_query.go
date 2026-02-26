package domain

import (
	"fmt"
	"strings"
)

const (
	techDefaultPage  = 1
	techDefaultLimit = 20
	techMaxLimit     = 200

	techMinTRL = 1
	techMaxTRL = 9
)

var techAllowedSortBy = map[string]struct{}{
	"name":            {},
	"trl":             {},
	"list_index":      {},
	"custom_metric_1": {},
	"custom_metric_2": {},
	"custom_metric_3": {},
	"custom_metric_4": {},
}

var techAllowedOrder = map[string]struct{}{
	"asc":  {},
	"desc": {},
}

// NormalizeAndValidateTechnologyListParams валидирует query-параметры списка технологий.
// ВАЖНО: здесь предполагается, что в твоём domain.TechnologyListParams поля называются:
// Page, Limit, SortBy, Order, TRLMin, TRLMax (инты, где 0 = не задано).
func NormalizeAndValidateTechnologyListParams(p *TechnologyListParams) error {
	// defaults
	if p.Page == 0 {
		p.Page = techDefaultPage
	}
	if p.Limit == 0 {
		p.Limit = techDefaultLimit
	}
	if strings.TrimSpace(p.SortBy) == "" {
		p.SortBy = "list_index"
	}
	if strings.TrimSpace(p.Order) == "" {
		p.Order = "asc"
	}

	// strict page/limit
	if p.Page < 1 {
		return fmt.Errorf("%w: page must be >= 1", ErrInvalid)
	}
	if p.Limit < 1 || p.Limit > techMaxLimit {
		return fmt.Errorf("%w: limit must be between 1 and %d", ErrInvalid, techMaxLimit)
	}

	// strict TRL (если заданы)
	if p.TRLMin != 0 && (p.TRLMin < techMinTRL || p.TRLMin > techMaxTRL) {
		return fmt.Errorf("%w: trl_min must be %d..%d", ErrInvalid, techMinTRL, techMaxTRL)
	}
	if p.TRLMax != 0 && (p.TRLMax < techMinTRL || p.TRLMax > techMaxTRL) {
		return fmt.Errorf("%w: trl_max must be %d..%d", ErrInvalid, techMinTRL, techMaxTRL)
	}
	if p.TRLMin != 0 && p.TRLMax != 0 && p.TRLMin > p.TRLMax {
		return fmt.Errorf("%w: trl_min must be <= trl_max", ErrInvalid)
	}

	// strict sort/order
	p.SortBy = strings.TrimSpace(p.SortBy)
	if _, ok := techAllowedSortBy[p.SortBy]; !ok {
		return fmt.Errorf("%w: sort_by must be one of name, trl, list_index, custom_metric_1..custom_metric_4", ErrInvalid)
	}

	p.Order = strings.ToLower(strings.TrimSpace(p.Order))
	if _, ok := techAllowedOrder[p.Order]; !ok {
		return fmt.Errorf("%w: order must be asc|desc", ErrInvalid)
	}

	return nil
}
