package prompt

func EstimateTokens(s string) int {
	if s == "" {
		return 0
	}
	runes := len([]rune(s))
	return int(float64(runes)/3.5*1.15) + 1
}
