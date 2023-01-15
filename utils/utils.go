package utils

import (
	"strings"

	humanize "github.com/dustin/go-humanize"
)

func FormatRupiah(amount int64) string {
	amountFloat64 := float64(amount)
	humanizeValue := humanize.CommafWithDigits(amountFloat64, 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	return "Rp " + stringValue
}
