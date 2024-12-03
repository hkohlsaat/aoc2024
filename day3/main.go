package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var filename = flag.String("input", "input.txt", "input for this assignment")

func main() {
	flag.Parse()

	b, err := os.ReadFile(*filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", *filename, err))
	}

	muls := FindMuls(string(b))
	result := Compute(muls)

	fmt.Printf("sum of multiplications: %d\n", result)

	muls = FindMuls(FilterInstructions(string(b)))
	result = Compute(muls)

	fmt.Printf("filtered instructions' result: %d\n", result)

}

func FilterInstructions(s string) string {
	filteredInstructions := []string{}
	for len(s) > 0 {
		var span string
		span, s = FindUntilDont(s)
		filteredInstructions = append(filteredInstructions, span)
		_, s = FindUntilDo(s)
	}
	return strings.Join(filteredInstructions, "")
}

func FindUntilDont(s string) (span string, rest string) {
	r := regexp.MustCompile(`don't\(\)`)
	loc := r.FindStringIndex(s)
	if loc == nil {
		return s, ""
	}
	return s[:loc[0]], s[loc[1]:]
}

func FindUntilDo(s string) (span string, rest string) {
	r := regexp.MustCompile(`do\(\)`)
	loc := r.FindStringIndex(s)
	if loc == nil {
		return s, ""
	}
	return s[:loc[0]], s[loc[1]:]
}

type Mul struct {
	X, Y int
}

func FindMuls(s string) []Mul {
	mulRegex := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	numberRegex := regexp.MustCompile(`\d{1,3}`)

	mulStrs := mulRegex.FindAllString(s, -1)
	muls := make([]Mul, 0, len(mulStrs))

	for _, mulStr := range mulStrs {
		numberStrs := numberRegex.FindAllString(mulStr, -1)
		x, _ := strconv.Atoi(numberStrs[0])
		y, _ := strconv.Atoi(numberStrs[1])
		muls = append(muls, Mul{x, y})
	}

	return muls
}

func Compute(muls []Mul) int {
	result := 0
	for _, mul := range muls {
		result += mul.X * mul.Y
	}
	return result
}
