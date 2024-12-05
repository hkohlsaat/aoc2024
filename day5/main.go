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

	rules, updates := ReadInput(strings.TrimSpace(string(b)))
	ordering := FromRules(rules)
	previouslyOrdered, newlyOrdered := ordering.Ordered(updates)

	fmt.Printf("sum of previously ordered middle pages: %d\n", SumMiddlePages(previouslyOrdered))
	fmt.Printf("sum of newly ordered middle pages: %d\n", SumMiddlePages(newlyOrdered))
}

type Pair struct {
	l, f int
}

func ReadInput(s string) ([]Pair, [][]int) {
	lines := strings.Split(s, "\n")

	rules := make([]Pair, 0, len(lines))

	readLines := 0
	for i, l := range lines {
		readLines++
		if len(l) == 0 {
			break
		}

		pair := strings.Split(l, "|")
		if len(pair) != 2 {
			panic(fmt.Sprintf("expected pair to have length of 2, has %d", len(pair)))
		}

		lead, err := strconv.Atoi(pair[0])
		if err != nil {
			panic(fmt.Sprintf("expected number, got %s: %s", pair[0], err))
		}
		follow, err := strconv.Atoi(pair[1])
		if err != nil {
			panic(fmt.Sprintf("expected number, got %s: %s", pair[0], err))
		}

		rules = append(rules, Pair{lead, follow})
		i++
	}

	updates := make([][]int, len(lines[readLines:]))
	for j, l := range lines[readLines:] {
		pageStrs := strings.Split(l, ",")
		pages := make([]int, len(pageStrs))
		for i, pageStr := range pageStrs {
			p, err := strconv.Atoi(pageStr)
			if err != nil {
				panic(fmt.Sprintf("error converting pages: %s", err))
			}
			pages[i] = p
		}
		updates[j] = pages
	}

	return rules, updates
}

func SumMiddlePages(updates [][]int) int {
	sum := 0
	for _, update := range updates {
		middleIndex := (len(update) / 2)
		middlePage := update[middleIndex]
		sum += middlePage
	}
	return sum
}

type OrderingRules struct {
	followers map[int][]int
	leaders   map[int][]int
}

func FromRules(rules []Pair) (or OrderingRules) {
	or.followers = make(map[int][]int)
	or.leaders = make(map[int][]int)
	for _, rule := range rules {
		or.followers[rule.l] = append(or.followers[rule.l], rule.f)
		or.leaders[rule.f] = append(or.leaders[rule.f], rule.l)
	}
	return or
}

func (rules OrderingRules) CanFollow(lead, follow int) bool {
	followers, ok := rules.followers[lead]
	if !ok {
		return false
	}
	for _, f := range followers {
		if f == follow {
			return true
		}
	}
	return false
}

func (rules OrderingRules) CanFollowAll(numbers []int) bool {
	for i := 1; i < len(numbers); i++ {
		if !rules.CanFollow(numbers[i-1], numbers[i]) {
			return false
		}
	}
	return true
}

func (rules OrderingRules) Order(pages []int) []int {
	if len(pages) == 0 {
		return []int{}
	}

	for i, page := range pages {
		pagesWithoutPage := append(append([]int(nil), pages[:i]...), pages[i+1:]...)
		order := rules.FindMostFollowers(page, pagesWithoutPage)
		if len(order)+1 == len(pages) {
			return append([]int{page}, order...)
		}
	}

	panic(fmt.Sprintf("no order for %v", pages))

	// followers := make([]int, 1, len(pages))
	// followers[0] = pages[0]
	// leaders := make([]int, 1, len(pages))
	// leaders[0] = pages[0]

	// for _, page := range pages[1:] {
	// 	if rules.CanFollow(followers[len(followers)-1], page) {
	// 		followers = append(followers, page)
	// 		continue
	// 	}
	// 	if rules.CanFollow(page, leaders[len(leaders)-1]) {
	// 		leaders = append(leaders, page)
	// 		continue
	// 	}
	// 	panic(fmt.Sprintf("no leader or follower found for %d", page))
	// }

	// reversedLeaders := make([]int, 0, len(leaders))
	// for i := len(leaders) - 1; i >= 0; i-- {
	// 	reversedLeaders = append(reversedLeaders, leaders[i])
	// }

	// return append(reversedLeaders, followers[1:]...)
}

func (rules OrderingRules) FindMostFollowers(page int, candidates []int) []int {
	if len(candidates) == 0 {
		return []int{}
	}
	if len(candidates) == 1 {
		if rules.CanFollow(page, candidates[0]) {
			return candidates
		}
	}

	followers := make([]int, 0, len(candidates))
	for i, candidate := range candidates {
		if rules.CanFollow(page, candidate) {
			c := append(append([]int(nil), candidates[:i]...), candidates[i+1:]...)
			ff := rules.FindMostFollowers(candidate, c)
			if len(ff)+1 > len(followers) {
				followers = append([]int{candidate}, ff...)
			}
		}
	}

	return followers
}

func (rules OrderingRules) OrderAll(unordered [][]int) [][]int {
	ordered := make([][]int, 0, len(unordered))
	for _, uo := range unordered {
		ordered = append(ordered, rules.Order(uo))
	}
	return ordered
}

func (rules OrderingRules) Ordered(numbers [][]int) ([][]int, [][]int) {
	ordered := make([][]int, 0, len(numbers))
	unordered := make([][]int, 0, len(numbers))
	for i := 0; i < len(numbers); i++ {
		if rules.CanFollowAll(numbers[i]) {
			ordered = append(ordered, numbers[i])
		} else {
			unordered = append(unordered, numbers[i])
		}
	}

	return ordered, rules.OrderAll(unordered)
}
