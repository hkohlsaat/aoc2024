package main

import "testing"

func TestCheckSafety(t *testing.T) {
	tests := []struct {
		report []int
		expect bool
	}{
		{[]int{1, 2}, true},
		{[]int{2, 2}, false},
		{[]int{1, 5}, false},
		{[]int{5, 1}, false},
		{[]int{5, 3}, true},
		{[]int{7, 6, 4, 2, 1}, true},
		{[]int{1, 2, 7, 8, 9}, false},
		{[]int{9, 7, 6, 2, 1}, false},
		{[]int{1, 3, 2, 4, 5}, false},
		{[]int{8, 6, 4, 4, 1}, false},
		{[]int{1, 3, 6, 7, 9}, true},
	}

	for _, test := range tests {
		got := CheckSafety(test.report)
		if got != test.expect {
			t.Errorf("%v: expected %t, got %t", test.report, test.expect, got)
		}
	}
}

func TestCheckDampenedSafety(t *testing.T) {
	tests := []struct {
		report []int
		expect bool
	}{
		{[]int{1, 3, 4, 4, 5}, true},
		{[]int{7, 6, 4, 2, 1}, true},
		{[]int{1, 2, 7, 8, 9}, false},
		{[]int{9, 7, 6, 2, 1}, false},
		{[]int{1, 3, 2, 4, 5}, true},
		{[]int{8, 6, 4, 4, 1}, true},
		{[]int{1, 3, 6, 7, 9}, true},
	}

	for _, test := range tests {
		got := CheckDampenedSafety(test.report)
		if got != test.expect {
			t.Errorf("%v: expected %t, got %t", test.report, test.expect, got)
		}
	}
}

func TestSplitReports(t *testing.T) {
	tests := []struct {
		raw    string
		expect [][]int
	}{
		{"3 4 5", [][]int{{3, 4, 5}}},
	}

testsLoop:
	for _, test := range tests {
		reports, err := SplitReports(test.raw)
		if err != nil {
			t.Errorf("returned error: %s", err)
			continue
		}

		if len(reports) != len(test.expect) {
			t.Errorf("report is not as expected: lengths do not match: expected %d, got %d", len(test.expect), len(reports))
			continue
		}

		for i := range reports {
			report := reports[i]
			expect := test.expect[i]

			if len(report) != len(expect) {
				t.Errorf("expected %v, got %v", expect, report)
				continue testsLoop
			}

			for j := range report {
				if report[j] != expect[j] {
					t.Errorf("expected %v, got %v", expect, report)
					continue testsLoop
				}
			}
		}
	}
}
