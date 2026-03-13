package service

import (
	"fmt"
	"math"
)

// buildTrendPosAndSegWidth: trendID -> position, segWidth = 2π / N
func buildTrendPosAndSegWidth(trendIDs []string) (map[string]int, float64) {
	trendPos := map[string]int{}
	for i, id := range trendIDs {
		trendPos[id] = i
	}

	segWidth := 2 * math.Pi
	if len(trendIDs) > 0 {
		segWidth = (2 * math.Pi) / float64(len(trendIDs))
	}
	return trendPos, segWidth
}

// radius: TRL 1..9 => [0..1]
func computeRadius(trl int) float64 {
	r := float64(trl-1) / 8.0
	if r < 0 {
		r = 0
	}
	if r > 1 {
		r = 1
	}
	return r
}

// angle: равномерно внутри сегмента тренда (через стабильный hash slug)
func computeAngle(trendPos map[string]int, segWidth float64, trendID, slug string) (float64, error) {
	pos, ok := trendPos[trendID]
	if !ok {
		return 0, fmt.Errorf("unknown trend id: %s", trendID)
	}
	u := hashUnit(slug)
	return float64(pos)*segWidth + u*segWidth, nil
}
