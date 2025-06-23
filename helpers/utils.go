package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RandFileName(ext string) string {
	filename := fmt.Sprintf("%s_%v", strings.ReplaceAll(uuid.NewString()[:6], "-", ""), time.Now().UnixMilli())
	if ext != "" {
		filename += fmt.Sprintf(".%s", ext)
	}
	return filename
}

func BinaryConverter(number int, bits int) []int {
	factor := number
	result := make([]int, bits)

	for factor >= 0 && number > 0 {
		factor = number % 2
		number /= 2
		result[bits-1] = factor
		bits--
	}
	return result
}
