package main

import (
	"bytes"
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

	sheet := SheetFromString(strings.TrimSpace(string(b)))
	possible := FindPossibleCalculations(sheet)
	sum := 0
	for _, p := range possible {
		sum += sheet.Calibrations[p].Result
	}
	fmt.Printf("Sum of possible calibrations: %d\n", sum)
}

type CalibrationSheet struct {
	Calibrations []*Calibration
}

type Calibration struct {
	Result   int
	Operands []int
}

func SheetFromString(s string) *CalibrationSheet {
	lines := strings.Split(s, "\n")
	calibrations := make([]*Calibration, len(lines))
	for i, line := range lines {
		calibrations[i] = CalibrationFromString(line)
	}

	return &CalibrationSheet{calibrations}
}

func CalibrationFromString(s string) *Calibration {
	parts := strings.Split(s, ": ")
	if len(parts) != 2 {
		panic("Calibration line does not have 2 parts")
	}
	strs := strings.Split(parts[1], " ")

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	operands := make([]int, len(strs))
	for i, str := range strs {
		operands[i], err = strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
	}

	return &Calibration{result, operands}
}

func (sheet *CalibrationSheet) String() string {
	b := bytes.Buffer{}
	for i, cali := range sheet.Calibrations {
		if i > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(cali.String())
	}
	return b.String()
}

func (cali *Calibration) String() string {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%d:", cali.Result))
	for _, op := range cali.Operands {
		b.WriteString(fmt.Sprintf(" %d", op))
	}
	return b.String()
}

type Operator string

const (
	Add         Operator = "+"
	Multiply    Operator = "*"
	Concatenate Operator = "||"
)

func FindPossibleCalculations(sheet *CalibrationSheet) []int {
	possible := make([]int, 0, len(sheet.Calibrations))
	for i, cal := range sheet.Calibrations {
		if FindOperators(cal) != nil {
			possible = append(possible, i)
		}
	}
	return possible
}

func FindOperators(cali *Calibration) []Operator {
	for operators := range IterOp(len(cali.Operands) - 1) {
		if Check(cali.Result, cali.Operands, operators) {
			return operators
		}
	}
	return nil
}

func Check(target int, operands []int, operators []Operator) bool {
	if len(operands) == 0 || len(operands)-1 != len(operators) {
		panic("cannot check calculation")
	}
	result := operands[0]
	for i, operand := range operands[1:] {
		operator := operators[i]
		switch operator {
		case Add:
			result = result + operand
		case Multiply:
			result = result * operand
		case Concatenate:
			for _, n := range strconv.Itoa(operand) {
				a, _ := strconv.Atoi(string(n))
				result = result*10 + a
			}
		default:
			panic("unknown operator")
		}
	}
	return result == target
}

func IterOp(total int) func(func([]Operator) bool) {
	return func(f func([]Operator) bool) {
		ops := make([]Operator, total)
		for expensiveOperations := range total + 1 {
			for concatenations := range expensiveOperations + 1 {
				if _continue := FillRecursively(ops, concatenations, expensiveOperations-concatenations, func() bool { return f(ops) }); !_continue {
					return
				}
			}
		}
	}
}

func FillRecursively(operators []Operator, concatenations, multiplications int, yield func() bool) bool {
	if len(operators) < concatenations+multiplications {
		return true
	}
	if len(operators) == 0 {
		return yield()
	}

	operators[0] = Add
	if _continue := FillRecursively(operators[1:], concatenations, multiplications, yield); !_continue {
		return false
	}

	if multiplications > 0 {
		operators[0] = Multiply
		if _continue := FillRecursively(operators[1:], concatenations, multiplications-1, yield); !_continue {
			return false
		}
	}

	if concatenations > 0 {
		operators[0] = Concatenate
		return FillRecursively(operators[1:], concatenations-1, multiplications, yield)
	}

	return true
}
