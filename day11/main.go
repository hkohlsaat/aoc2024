package main

import (
	"flag"
	"fmt"
	"os"
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
	stones := StringToStones(strings.TrimSpace(string(b)))
	memory := make(StoneMemory)
	count25 := CountAfterNBlinks(memory, stones, 0, 25)
	fmt.Printf("Count after 25 blinks: %d\n", count25)

	count75 := CountAfterNBlinks(memory, stones, 0, 75)
	fmt.Printf("Count after 75 blinks: %d\n", count75)
}

func StringToStones(s string) []int {
	numberStrs := strings.Split(s, " ")
	stones := make([]int, 0, len(numberStrs))
	for _, numberStr := range numberStrs {
		stone, err := strconv.Atoi(numberStr)
		if err != nil {
			panic(err)
		}
		stones = append(stones, stone)
	}
	return stones
}

type StoneMemory = map[[2]int]int

func CountAfterNBlinks(memory StoneMemory, stones []int, blinks int, maxBlinks int) int {
	if blinks == maxBlinks {
		return len(stones)
	}

	count := 0
	for _, stone := range stones {
		thisCount := 0
		if c, ok := memory[[2]int{stone, maxBlinks - blinks}]; ok {
			thisCount = c
		} else if stone == 0 {
			thisCount = CountAfterNBlinks(memory, []int{1}, blinks+1, maxBlinks)
		} else if stoneStr := strconv.Itoa(stone); len(stoneStr)%2 == 0 {
			left, _ := strconv.Atoi(stoneStr[:len(stoneStr)/2])
			right, _ := strconv.Atoi(stoneStr[len(stoneStr)/2:])
			thisCount = CountAfterNBlinks(memory, []int{left, right}, blinks+1, maxBlinks)
		} else {
			thisCount = CountAfterNBlinks(memory, []int{stone * 2024}, blinks+1, maxBlinks)
		}
		memory[[2]int{stone, maxBlinks - blinks}] = thisCount
		count += thisCount
	}
	return count
}
