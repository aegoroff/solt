package fw

// Percent calculates percent value
func Percent(value int64, total int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(value) / float64(total)) * 100
}
