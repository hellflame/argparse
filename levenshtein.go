package argparse


func decideMatch(target string, candidates []string) string {
	ldArray := make([]int, len(candidates))
	for i, c := range candidates {
		ldArray[i] = levDistance(target, c)
	}
	match := min(ldArray...)
	if match >= len(target) {  // too many diff
		return ""
	}
	matchCandidates := make(map[int]string)
	var matchKeys []int
	for i, ld := range ldArray {
		if ld == match {
			wordL := len(candidates[i])
			matchKeys = append(matchKeys, wordL)
			matchCandidates[wordL] = candidates[i]
		}
	}
	shortest := min(matchKeys...)
	return matchCandidates[shortest]
}

func levDistance(a, b string) int {
	la := len(a)
	lb := len(b)
	matrix := make([][]int, la + 1)
	for i := range matrix {
		matrix[i] = make([]int, lb + 1)
	}
	for i := 0; i <= la; i += 1 {
		matrix[i][0] = i
	}
	for i := 0; i <= lb; i += 1 {
		matrix[0][i] = i
	}
	for i := 1; i <= la; i += 1 {
		for j := 1; j <= lb; j += 1 {
			cost := 1
			if a[i - 1] == b[j - 1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i - 1][j - 1] + cost,
				matrix[i][j - 1] + 1,
				matrix[i - 1][j] + 1)
		}
	}
	return matrix[la][lb]
}

func _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min(candidate ... int) int {
	if len(candidate) == 1 {
		return candidate[0]
	}
	if len(candidate) == 2 {
		return _min(candidate[0], candidate[1])
	}
	return _min(candidate[0], min(candidate[1:]...))
}
