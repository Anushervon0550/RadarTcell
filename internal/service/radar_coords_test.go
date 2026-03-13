package service

import "testing"

func TestBuildTrendPosAndSegWidth(t *testing.T) {
	ids := []string{"t1", "t2", "t3"}
	pos, seg := buildTrendPosAndSegWidth(ids)
	if seg <= 0 {
		t.Fatalf("expected segWidth > 0, got %v", seg)
	}
	if pos["t1"] != 0 || pos["t2"] != 1 || pos["t3"] != 2 {
		t.Fatalf("unexpected trend positions: %#v", pos)
	}
}

func TestComputeRadiusBounds(t *testing.T) {
	cases := []struct {
		trl  int
		want float64
	}{
		{1, 0},
		{9, 1},
		{0, 0},
		{10, 1},
	}
	for _, c := range cases {
		if got := computeRadius(c.trl); got != c.want {
			t.Fatalf("trl %d: expected %v, got %v", c.trl, c.want, got)
		}
	}
}

func TestComputeAngleRange(t *testing.T) {
	pos, seg := buildTrendPosAndSegWidth([]string{"t1", "t2"})
	angle1, err := computeAngle(pos, seg, "t1", "slug-a")
	if err != nil {
		t.Fatalf("unexpected err for t1: %v", err)
	}
	if angle1 < 0 || angle1 >= seg {
		t.Fatalf("angle for t1 out of range: %v", angle1)
	}
	angle2, err := computeAngle(pos, seg, "t2", "slug-b")
	if err != nil {
		t.Fatalf("unexpected err for t2: %v", err)
	}
	if angle2 < seg || angle2 >= 2*seg {
		t.Fatalf("angle for t2 out of range: %v", angle2)
	}
}

func TestComputeAngle_UnknownTrendIDReturnsError(t *testing.T) {
	pos, seg := buildTrendPosAndSegWidth([]string{"t1", "t2"})
	_, err := computeAngle(pos, seg, "missing-trend", "slug-a")
	if err == nil {
		t.Fatal("expected error for unknown trend id, got nil")
	}
}

