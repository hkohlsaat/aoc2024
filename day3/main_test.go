package main

import (
	"testing"
)

func TestFindMuls(t *testing.T) {
	tests := []struct {
		raw  string
		muls []Mul
	}{
		{"mul(1,2)", []Mul{{1, 2}}},
		{"xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))", []Mul{{2, 4}, {5, 5}, {11, 8}, {8, 5}}},
	}

testsLoop:
	for _, test := range tests {
		muls := FindMuls(test.raw)
		if len(muls) != len(test.muls) {
			t.Errorf("expected %d muls, got %d", len(test.muls), len(muls))
			continue
		}

		for i, mul := range muls {
			expect := test.muls[i]
			if mul.X != expect.X || mul.Y != expect.Y {
				t.Errorf("at index %d: expected mul(%d,%d), got mul(%d,%d)", i, expect.X, expect.Y, mul.X, mul.Y)
				continue testsLoop
			}
		}
	}
}

func TestComputeMuls(t *testing.T) {
	tests := []struct {
		muls   []Mul
		expect int
	}{
		{[]Mul{{1, 4}, {4, 6}, {2, 2}}, 4 + 24 + 4},
	}

	for _, test := range tests {
		got := Compute(test.muls)
		if got != test.expect {
			t.Errorf("expected %d, got %d", test.expect, got)
		}
	}
}

func TestFindUntilDont(t *testing.T) {
	tests := []struct {
		input      string
		span, rest string
	}{
		{"asdfdon't()fdsa", "asdf", "fdsa"},
		{"asdf", "asdf", ""},
		{"", "", ""},
	}

	for _, test := range tests {
		s, r := FindUntilDont(test.input)
		if s != test.span {
			t.Errorf("expected %s, got %s", test.span, s)
		} else if r != test.rest {
			t.Errorf("expected %s, got %s", test.rest, r)
		}
	}
}

func TestFindUntilDo(t *testing.T) {
	tests := []struct {
		input      string
		span, rest string
	}{
		{"asdfdo()fdsa", "asdf", "fdsa"},
		{"asdf", "asdf", ""},
		{"", "", ""},
	}

	for _, test := range tests {
		s, r := FindUntilDo(test.input)
		if s != test.span {
			t.Errorf("expected %s, got %s", test.span, s)
		} else if r != test.rest {
			t.Errorf("expected %s, got %s", test.rest, r)
		}
	}
}

func TestFilterInstructions(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"adon't()bcdo()defdon't()do()don't()ghdo()i", "adefi"},
	}

	for _, test := range tests {
		got := FilterInstructions(test.input)
		if got != test.expect {
			t.Errorf("expected %s, got %s", test.expect, got)
		}
	}
}
