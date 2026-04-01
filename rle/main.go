package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrDigitsInvalidFormat = errors.New("invalid string format: digit cannot precede a character")
var ErrEscapeInvalidFormat = errors.New("invalid string format: dangling escape symbol as the last character")

func UnpackRune(r []rune) (string, error) {
	var b strings.Builder
	b.Grow(len(r))

	leftPtr := 0
	var isEscaped bool
	for rightPtr, digRune := range r {
		if unicode.IsDigit(digRune) {

			// If first char is digit, it is incorrect format.
			if rightPtr == 0 {
				return "", ErrDigitsInvalidFormat
			}

			// This digit taken into account already.
			if unicode.IsDigit(r[rightPtr-1]) {
				continue
			}

			// Write chars do not require unpacking (just like abc).
			for rightPtr > leftPtr+1 {
				if (isEscaped && r[leftPtr] == '\\') || r[leftPtr] != '\\' {
					b.WriteRune(r[leftPtr])
				}
				if !isEscaped && r[leftPtr] == '\\' {
					isEscaped = true
				} else if isEscaped {
					isEscaped = false
				}
				leftPtr++
			}

			// Make a num if it is not just a digit (a10 -> aaaaaaaaaa).
			var digs strings.Builder
			// Escaping symbol matters (b\110 -> b1111111111; e\\3 -> e\\\).
			if !isEscaped && r[leftPtr] == '\\' {
				leftPtr++ // digit to be repeated, not escape.
			} else {
				digs.WriteRune(digRune)
			}
			// Count amount of a char to be decoded.
			for rightPtr+1 < len(r) && unicode.IsDigit(r[rightPtr+1]) {
				rightPtr++
				digs.WriteRune(r[rightPtr])
			}
			// Len is zero when it is just an escape with digit (just like a\1 -> a1).
			if digs.Len() == 0 {
				digs.WriteRune('1')
			}
			digInt, err := strconv.Atoi(digs.String())
			if err != nil {
				return "", fmt.Errorf("failed to convert %s to int: %w", digs.String(), err)
			}

			// Write char require unpacking (just like a3b2 -> aaabb).
			strToCat := strings.Repeat(string(r[leftPtr]), digInt)
			b.WriteString(strToCat)

			if rightPtr+1 < len(r) {
				leftPtr = rightPtr + 1
			}
			isEscaped = false
		}
	}

	// Write leftover chars do not require unpacking.
	if len(r) > 0 && !unicode.IsDigit(r[len(r)-1]) {
		for leftPtr < len(r) {
			if (isEscaped && r[leftPtr] == '\\') || r[leftPtr] != '\\' {
				b.WriteRune(r[leftPtr])
			}
			if !isEscaped && r[leftPtr] == '\\' {
				isEscaped = true
			} else if isEscaped {
				isEscaped = false
			}
			leftPtr++
		}
	}

	if isEscaped && r[len(r)-1] == '\\' {
		return "", ErrEscapeInvalidFormat
	}

	return b.String(), nil
}

func main() {
	var str string
	fmt.Print("enter str: ")
	fmt.Scan(&str)
	//str = "abc\\"
	fmt.Println("cur str:", str)

	modStr, err := UnpackRune([]rune(str))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("mod str:", modStr)
}
