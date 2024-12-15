package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var input = func() string {
	filename := "input_test.txt"
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("could not read file %s: %s\n", filename, err))
	}
	return string(b)
}()

func TestFromStringAndString(t *testing.T) {
	sheet := SheetFromString(input)
	got := sheet.String()
	if got != input {
		t.Errorf("expected:\n%s\n\ngot:\n%s", input, got)
	}
}

func TestFindPossibleCalibrations(t *testing.T) {
	tests := []struct {
		sheet  *CalibrationSheet
		expect []int
	}{
		{SheetFromString(input), []int{0, 1, 3, 4, 6, 8}},
	}

testLoop:
	for _, test := range tests {
		got := FindPossibleCalculations(test.sheet)
		if len(got) != len(test.expect) {
			t.Errorf("expected %v, got %v", test.expect, got)
			continue
		}
		for i, n := range got {
			if n != test.expect[i] {
				t.Errorf("expected %v, got %v", test.expect, got)
				continue testLoop
			}
		}
	}
}

func TestFindOperators(t *testing.T) {
	tests := []struct {
		cali *Calibration
		ops  []Operator
	}{
		{CalibrationFromString("3: 1 2"), []Operator{"+"}},
		{CalibrationFromString("6: 2 3"), []Operator{"*"}},
		{&Calibration{((3*4+5)*12+23)*3 + 3, []int{3, 4, 5, 12, 23, 3, 3}}, []Operator{"*", "+", "*", "+", "*", "+"}},
	}

testLoop:
	for tt, test := range tests {
		ops := FindOperators(test.cali)
		if len(ops) != len(test.ops) {
			t.Errorf("%d: expected %d operators, got %d", tt, len(test.ops), len(ops))
			break
		}

		for i, op := range ops {
			expect := test.ops[i]
			if op != expect {
				t.Errorf("%d: expected %s, got %s", tt, operatorsToString(test.ops), operatorsToString(ops))
				continue testLoop
			}
		}
	}
}

func operatorsToString(ops []Operator) string {
	b := bytes.Buffer{}
	for i, o := range ops {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(o))
	}
	return b.String()
}
