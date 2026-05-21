package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/k0kubun/pp"
)

func main() {
	strs := []string{"Пятак", "пяткА", "тЯпка", "листоК", "слитОк", "столик", "сТол"}
	AnagramSearch(os.Stdout, strs)
}

func AnagramSearch(out io.Writer, strs []string) {
	anagramSet := make(map[string][]string, len(strs))

	for _, str := range strs {
		strLowered := strings.ToLower(str)

		runes := []rune(strLowered)
		slices.Sort(runes)
		sortedStr := string(runes)

		anagramSet[sortedStr] = append(anagramSet[sortedStr], strLowered)
	}

	for _, anagramSlice := range anagramSet {
		if len(anagramSlice) > 1 {
			slices.Sort(anagramSlice)
			pp.Fprintf(out, "- \"%s\": %v\n", anagramSlice[0], fmt.Sprintf("%s", anagramSlice))
		}
	}
}
