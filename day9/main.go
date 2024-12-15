package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var filename = flag.String("input", "input.txt", "input for this assignment")

func main() {
	flag.Parse()

	b, err := os.ReadFile(*filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", *filename, err))
	}

	diskMap := DiscMapFromString(string(b))
	checksum := ComputeChecksum(diskMap)
	fmt.Printf("Checksum (single block reordering): %d\n", checksum)

	memorySpans := MemorySpansFromString(string(b))
	memorySpans = ReorderSpans(memorySpans)
	checksum = ComputeChecksumFromMemorySpans(memorySpans)
	fmt.Printf("Checksum (whole file reordering): %d\n", checksum)
}

func DiscMapFromString(s string) []int {
	discMap := make([]int, len(s))
	for i, r := range s {
		err := error(nil)
		discMap[i], err = strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
	}
	return discMap
}

func ComputeChecksum(diskMap []int) int {
	checksum := 0
	position := 0
	for id, length := range ReadFromFront(diskMap) {
		for range length {
			checksum += id * position
			position++
		}
	}
	return checksum
}

func ReadFromBack(diskMap []int) (id, length, read int) {
	if len(diskMap)%2 == 0 {
		diskMap = diskMap[:len(diskMap)-1]
		read = 1
	}

	id = len(diskMap) / 2
	length = diskMap[len(diskMap)-1]
	read += 1
	return id, length, read
}

func ReadFromFront(diskMap []int) func(yield func(id, length int) bool) {
	return func(yield func(id, length int) bool) {
		var (
			leftOverId     = 0
			leftOverLength = 0
		)

		for i := 0; i < len(diskMap); i++ {
			length := diskMap[i]
			if i%2 == 0 {
				if !yield(i/2, length) {
					return
				}
				continue
			}

			for length > 0 {
				var id, l int
				if leftOverLength > 0 {
					id = leftOverId
					l = leftOverLength
					leftOverLength = 0
				} else if i+1 < len(diskMap) {
					var r int
					id, l, r = ReadFromBack(diskMap[i+1:])
					diskMap = diskMap[:len(diskMap)-r]
					id += i/2 + 1
				} else {
					return
				}

				if l > length {
					leftOverId = id
					leftOverLength = l - length
					l = length
				}
				if !yield(id, l) {
					return
				}
				length -= l
			}
		}

		if leftOverLength > 0 {
			yield(leftOverId, leftOverLength)
		}
	}
}

type MemoryType int

const (
	Empty MemoryType = iota
	File
)

type MemorySpan struct {
	Type   MemoryType
	Length int
	FileId int
}

func MemorySpansFromString(s string) []MemorySpan {
	spans := make([]MemorySpan, 0, len(s))
	for i, r := range s {
		span := MemorySpan{}
		var err error
		if span.Length, err = strconv.Atoi(string(r)); err != nil {
			panic(err)
		}

		if i%2 == 0 {
			span.Type = File
			span.FileId = i / 2
		}

		spans = append(spans, span)
	}
	return spans
}

func MemorySpansToString(spans []MemorySpan) string {
	b := &bytes.Buffer{}
	for _, s := range spans {
		r := "."
		if s.Type == File {
			r = strconv.Itoa(s.FileId)
		}
		for range s.Length {
			b.WriteString(r)
		}
	}
	return b.String()
}

func ReorderSpans(spans []MemorySpan) []MemorySpan {
	lowestId := len(spans)
	for i := len(spans) - 1; i >= 0; i-- {
		span := spans[i]
		if span.Type == Empty || span.FileId >= lowestId {
			continue
		}

		freeSpaceIndex := -1
		freeSpaceSpan := MemorySpan{}
		for j, _span := range spans[:i+1] {
			if _span.Type == Empty && _span.Length >= span.Length {
				freeSpaceIndex = j
				freeSpaceSpan = _span
				break
			}
		}

		if freeSpaceIndex < 0 {
			continue
		}

		var (
			beforeFreeSpaceIndex = spans[:freeSpaceIndex]
			insert               = []MemorySpan{spans[i]}
			middlePart           = spans[freeSpaceIndex+1 : i]
			spanReplacement      = []MemorySpan{{Length: span.Length}}
			afterSpan            = spans[i+1:]
		)

		if freeSpaceSpan.Length > span.Length {
			insert = append(insert, MemorySpan{Length: freeSpaceSpan.Length - span.Length})
			i += 1
		}

		spans = append(append(append(append(append([]MemorySpan(nil),
			beforeFreeSpaceIndex...),
			insert...),
			middlePart...),
			spanReplacement...),
			afterSpan...)

		lowestId = span.FileId
	}
	return spans
}

func ComputeChecksumFromMemorySpans(spans []MemorySpan) int {
	checksum := 0
	position := 0
	for _, span := range spans {
		if span.Type == Empty {
			position += span.Length
			continue
		}

		for range span.Length {
			checksum += position * span.FileId
			position++
		}
	}
	return checksum
}
