package utils

import (
	"strings"
	"unicode"

	humanize "github.com/dustin/go-humanize"
)

func FormatRupiah(amount int64) string {
	amountFloat64 := float64(amount)
	humanizeValue := humanize.CommafWithDigits(amountFloat64, 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	return "Rp " + stringValue
}

func CapitalizeFirstChar(s string) string {
	runeSlice := []rune(s)

	if len(runeSlice) > 0 && unicode.IsLetter(runeSlice[0]) {
		runeSlice[0] = unicode.ToUpper(runeSlice[0])
	}

	return string(runeSlice)
}
