package main

import (
	"testing"
)

func TestSumMiddlePages(t *testing.T) {
	tests := []struct {
		updates [][]int
		expect  int
	}{
		{[][]int{{1, 3, 5}, {2, 4, 6}}, 3 + 4},
	}

	for _, test := range tests {
		got := SumMiddlePages(test.updates)
		if got != test.expect {
			t.Errorf("expected %d, got %d", test.expect, got)
		}
	}
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		input   string
		rules   []Pair
		updates [][]int
	}{
		{"1|2", []Pair{{1, 2}}, nil},
		{"1|2\n2|3\n\n1,2,3", []Pair{{1, 2}, {2, 3}}, [][]int{{1, 2, 3}}},
	}

testLoop:
	for _, test := range tests {
		rules, updates := ReadInput(test.input)
		if len(rules) != len(test.rules) {
			t.Errorf("expected %d rules, got %d rules", len(test.rules), len(rules))
			continue
		}
		if len(updates) != len(test.updates) {
			t.Errorf("expected %d updates, got %d updates", len(test.updates), len(updates))
			continue
		}

		for i, rule := range rules {
			expect := test.rules[i]
			if rule.l != expect.l || rule.f != expect.f {
				t.Errorf("expected %v, got %v", expect, rule)
				continue testLoop
			}
		}

		for i, update := range updates {
			expect := test.updates[i]
			if len(update) != len(expect) {
				t.Errorf("expected %d pages, got %d", len(expect), len(update))
				continue testLoop
			}
			for j, page := range update {
				if page != expect[j] {
					t.Errorf("expected page %d, got %d", expect[j], page)
				}
			}
		}
	}
}

func TestFromRules(t *testing.T) {
	tests := []struct {
		rules       []Pair
		falseChecks []Pair
	}{
		{[]Pair{{1, 2}, {3, 4}}, []Pair{{2, 3}}},
	}

	for _, test := range tests {
		r := FromRules(test.rules)
		for _, rule := range test.rules {
			if !r.CanFollow(rule.l, rule.f) {
				t.Errorf("expected to allow %v, but is not allowed", rule)
			}
		}
		for _, falseCheck := range test.falseChecks {
			if r.CanFollow(falseCheck.l, falseCheck.f) {
				t.Errorf("expected to allow %v, but is not allowed", falseCheck)
			}
		}
	}
}

func TestCanFollow(t *testing.T) {
	rules := OrderingRules{
		followers: map[int][]int{
			9: {1},
		},
	}
	tests := []struct {
		lead, follow int
		expect       bool
	}{
		{1, 9, false},
		{2, 3, false},
		{9, 1, true},
		{9, 4, false},
	}

	for _, test := range tests {
		got := rules.CanFollow(test.lead, test.follow)
		if got != test.expect {
			t.Errorf("expected %t, got %t", test.expect, got)
		}
	}
}

func TestOrder(t *testing.T) {
	rules := FromRules([]Pair{{1, 2}, {4, 3}, {3, 1}})
	tests := []struct {
		pages  []int
		expect []int
	}{
		{[]int{1, 2, 3, 4}, []int{4, 3, 1, 2}},
	}

	for _, test := range tests {
		got := rules.Order(test.pages)
		if len(got) != len(test.expect) {
			t.Errorf("expected %d ordered pages, got %d", len(test.expect), len(got))
			continue
		}
		if !intSlicesSame(got, test.expect) {
			t.Errorf("expected %v, got %v", test.expect, got)
		}
	}
}

func intSlicesSame(aSlice, bSlice []int) bool {
	for i, a := range aSlice {
		if a != bSlice[i] {
			return false
		}
	}
	return true
}

func TestOrder2(t *testing.T) {
	r, _ := ReadInput(`47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`)
	rules := FromRules(r)
	tests := []struct {
		pages  []int
		expect []int
	}{
		{[]int{75, 97, 47, 61, 53}, []int{97, 75, 47, 61, 53}},
		{[]int{61, 13, 29}, []int{61, 29, 13}},
		{[]int{97, 13, 75, 29, 47}, []int{97, 75, 47, 29, 13}},
	}

	for _, test := range tests {
		got := rules.Order(test.pages)
		if len(got) != len(test.expect) {
			t.Errorf("expected %d ordered pages, got %d, %v", len(test.expect), len(got), got)
			continue
		}
		if !intSlicesSame(got, test.expect) {
			t.Errorf("expected %v, got %v", test.expect, got)
		}
	}
}
