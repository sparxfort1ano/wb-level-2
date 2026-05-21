package options

import (
	"cmp"
	"regexp"
	"strconv"
	"strings"
)

func splitByTabAndGetNElem(str string, n int) string {
	currTokens := strings.Split(str, "\t")
	index := n - 1
	if index >= 0 && index < len(currTokens) {
		return currTokens[index]
	}
	return ""
}

func getMonthValue(str string) int {
	var monthValues = map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
		"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8,
		"SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
	}

	strStripped := strings.TrimLeft(str, " \t")
	if len([]rune(strStripped)) < 3 {
		return 0
	}

	prefix := strings.ToUpper(strStripped[:3])
	return monthValues[prefix]
}

var reSize = regexp.MustCompile(`^([0-9\.\-\+]*)(.)?`)

func getSizeValue(str string) float64 {
	var sizeValues = map[string]int{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": 1024 * 1024 * 1024 * 1024,
	}

	strStripped := strings.TrimLeft(str, " \t")
	matches := reSize.FindStringSubmatch(strStripped)

	valFloat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		valFloat = 0
	}

	suffix := matches[2]
	valSuffix, ok := sizeValues[suffix]
	if !ok {
		valSuffix = 1
	}

	return valFloat * float64(valSuffix)
}

var reFloat = regexp.MustCompile(`^[+-]?\d+(\.\d+)?`)

func getNumericValue(str string) float64 {
	strStripped := strings.TrimLeft(str, " \t")
	match := reFloat.FindString(strStripped)

	numFloat, err := strconv.ParseFloat(match, 64)
	if err != nil {
		numFloat = 0
	}

	return numFloat
}

// Compare compares two adjacent strings according to the sorting options.
// Returns -1 if left is less than y, 1 if right is less than x, 0 if there are equal.
func (o *Options) Compare(left, right string) int {
	var result int

	if o.ByColumn > 0 {
		left = splitByTabAndGetNElem(left, o.ByColumn)
		right = splitByTabAndGetNElem(right, o.ByColumn)
	}

	if o.IgnoreTrailingBlanks {
		left = strings.TrimRight(left, " \t\n\r")
		right = strings.TrimRight(right, " \t\n\r")
	}

	if o.ByMonth {
		result = cmp.Compare(getMonthValue(left), getMonthValue(right))
	} else if o.BySize {
		result = cmp.Compare(getSizeValue(left), getSizeValue(right))
	} else if o.ByValue {
		result = cmp.Compare(getNumericValue(left), getNumericValue(right))
	}

	if result == 0 {
		result = cmp.Compare(left, right)
	}

	if o.Reverse {
		result *= (-1)
	}

	return result
}

// Equal checks whether adjacent elements are equal.
func (o *Options) Equal(left, right string) bool {
	return o.Compare(left, right) == 0
}
