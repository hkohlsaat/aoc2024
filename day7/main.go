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

type Operator rune

const (
	Add = '+'
	Mul = '*'
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
	// return RecursiveOperators(cali.Result, cali.Operands...)
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
		case Mul:
			result = result * operand
		default:
			panic("unknown operator")
		}
	}
	return result == target
}

func RecursiveOperators(target int, operands ...int) []Operator {
	if len(operands) == 1 {
		if target == operands[0] {
			return []Operator{}
		} else {
			return nil
		}
	}

	l, r := operands[0], operands[1]
	for _, operator := range []Operator{Add, Mul} {
		intermediateResult := 0
		switch operator {
		case Add:
			intermediateResult = l + r
		case Mul:
			intermediateResult = l * r
		default:
			panic("unknwon case")
		}
		if intermediateResult > target {
			return nil
		}
		newOperands := append([]int{intermediateResult}, operands[2:]...)
		ops := RecursiveOperators(target, newOperands...)
		if ops != nil {
			return append([]Operator{operator}, ops...)
		}
	}
	return nil
}

func IterOp(total int) func(func([]Operator) bool) {
	return func(f func([]Operator) bool) {
		ops := make([]Operator, total)
		for i := range total {
			cont := DistributeMuls(ops, i, func() bool { return f(ops) })
			if !cont {
				return
			}
		}
	}
}

func DistributeMuls(operators []Operator, multiplications int, yield func() bool) (cont bool) {
	if len(operators) < multiplications {
		return true
	}
	if len(operators) == 0 {
		return yield()
	}
	for _, occupy := range []bool{false, true} {
		if occupy {
			operators[0] = Mul
			cont = DistributeMuls(operators[1:], multiplications-1, yield)
		} else {
			operators[0] = Add
			cont = DistributeMuls(operators[1:], multiplications, yield)
		}
		if !cont {
			return false
		}
	}
	return true
}
