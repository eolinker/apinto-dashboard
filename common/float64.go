package common

import (
	"fmt"
	"strconv"
)

// FormatFloat64 float64保留places位小数
func FormatFloat64(f float64, places int) float64 {
	formatStr := "%." + strconv.Itoa(places) + "f"
	f64, _ := strconv.ParseFloat(fmt.Sprintf(formatStr, f), 64)
	return f64
}
