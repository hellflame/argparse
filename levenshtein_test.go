package argparse

import (
	"strings"
	"testing"
)

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
	if decideMatch("linux", []string{"linux", "a", "b"})[0] != "linux" {
		t.Error("failed to match same word")
		return
	}
	if decideMatch("linux", []string{"linu", "a", "b"})[0] != "linu" {
		t.Error("failed to match")
		return
	}
	if decideMatch("linux", []string{"linuxa", "linua", "b"})[0] != "linua" {
		t.Error("failed to choose shorter one")
		return
	}
	if decideMatch("linux", []string{"bilibili", "ok", "z"})[0] != "" {
		t.Error("failed to stop match")
		return
	}

	if strings.Join(decideMatch("linux", []string{"linua", "linub", "iinux"}), ",") != "linua,linub,iinux" {
		t.Error("failed to match multi")
		return
	}
}
