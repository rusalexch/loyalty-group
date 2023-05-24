package utils

import "math"

func ToCent(price float64) uint {
	return uint(math.Round(price * 100))
}

func FromCent(price uint) float64 {
	return float64(price) / 100
}
