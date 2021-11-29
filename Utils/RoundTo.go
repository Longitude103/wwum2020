package Utils

import "math"

// RoundTo rounds a float number to a specified number of decimal places.
func RoundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}
