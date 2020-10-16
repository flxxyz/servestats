package utils

import (
	"fmt"
	"math"
)

func LogN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func HumanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(LogN(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.2f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// FileSize 格式化文件大小
func FileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return HumanateBytes(s, 1024, sizes)
}
