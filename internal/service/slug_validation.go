package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateSlugValue(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("%w: slug is required", domain.ErrInvalid)
	}
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("%w: slug must match ^[a-z0-9]+(?:-[a-z0-9]+)*$", domain.ErrInvalid)
	}
	return nil
}
