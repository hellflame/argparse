package argparse

import "testing"

func TestLevDistance(t *testing.T) {
	for _, row :=  range [][]string{{"a", "a"}, {"well", "well"}, {"-linux", "-linux"}} {
		if levDistance(row[0], row[1]) != 0 {
			t.Error("should match")
			return
		}
	}

	for _, row :=  range [][]string{{"a", "b"}, {"well", "xell"}, {"linux", "-linux"}} {
		if levDistance(row[0], row[1]) != 1 {
			t.Error("should be only 1")
			return
		}
	}
}


func TestLevDecide(t *testing.T) {
	if decideMatch("linux", []string{"linux", "a", "b"}) != "linux" {
		t.Error("failed to match same word")
		return
	}
	if decideMatch("linux", []string{"linu", "a", "b"}) != "linu" {
		t.Error("failed to match")
		return
	}
}
