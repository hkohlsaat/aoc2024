package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var filename = flag.String("input", "input.txt", "input for this assignment")

func main() {
	flag.Parse()

	b, err := os.ReadFile(*filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", *filename, err))
	}
	input := strings.TrimSpace(string(b))

	_ = input
}
