package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/k0kubun/pp"
)

type AnagramCollection struct {
	FirstStr string
	Strs     []string
}

func NewAnagramCollection(firstStr string, strs []string) *AnagramCollection {
	return &AnagramCollection{
		FirstStr: firstStr,
		Strs:     strs,
	}
}

func main() {
	strs := []string{"Пятак", "пяткА", "тЯпка", "слитОк", "листоК", "столик", "сТол"}
	AnagramSearch(os.Stdout, strs)
}

func AnagramSearch(out io.Writer, strs []string) {
	anagramSet := make(map[string]*AnagramCollection, len(strs))

	for _, str := range strs {
		strLowered := strings.ToLower(str)

		runes := []rune(strLowered)
		slices.Sort(runes)
		sortedStr := string(runes)

		if _, ok := anagramSet[sortedStr]; !ok {
			anagramSet[sortedStr] = NewAnagramCollection(strLowered, make([]string, 0, len(str)))
		}

		anagramSet[sortedStr].Strs = append(anagramSet[sortedStr].Strs, strLowered)
	}

	for _, anagramCollection := range anagramSet {
		if len(anagramCollection.Strs) > 1 {
			slices.Sort(anagramCollection.Strs)
			pp.Fprintf(
				out,
				"- \"%s\": %v\n",
				anagramCollection.FirstStr,
				fmt.Sprintf("%s", anagramCollection.Strs))
		}
	}
}
