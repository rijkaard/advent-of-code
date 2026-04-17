package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MathEntryType = int

const (
	MathEntryNumber MathEntryType = iota
	MathEntryOperation
	MathEntryNone
)

type OperationType = int

const (
	OperationPlus OperationType = iota
	OperationTimes
	OperationNone
)

// func yieldFromLine(yield func(MathEntryType, uint, OperationType) bool) {
func yieldEntriesFromLine(line string) <-chan string {
	out := make(chan string)
	go func() {
		parts := strings.Split(line, " ")
		for idx := range parts {
			if parts[idx] != "" {
				out <- parts[idx]
			}
		}
		close(out)
	}()
	return out
}

func parseEntry(entry string) (MathEntryType, uint, OperationType, error) {
	num, err := strconv.Atoi(entry)
	if err == nil {
		return MathEntryNumber, uint(num), OperationNone, nil
	}
	if entry == "*" {
		return MathEntryOperation, 0, OperationTimes, nil
	} else if entry == "+" {
		return MathEntryOperation, 0, OperationPlus, nil
	}
	return MathEntryNone, 0, OperationNone, fmt.Errorf("invalid entry: %s", entry)
}

func performPlus(wot []uint) (out uint) {
	out = 0
	for x := range wot {
		out += wot[x]
	}
	return out
}

func performTimes(wot []uint) (out uint) {
	out = 1
	for x := range wot {
		out *= wot[x]
	}
	return out
}

func aoc6a() uint {
	problems := make([][]uint, 0)
	password := uint(0)
	for line := range parseLinesFromStdin {
		cur_entry_type := MathEntryNone
		cur_idx := 0
		for entry := range yieldEntriesFromLine(line) {
			entry_type, entry_number, entry_op, err := parseEntry(entry)
			if err != nil {
				fmt.Println(err)
				return 0
			}
			if cur_entry_type != MathEntryNone && entry_type != cur_entry_type {
				fmt.Printf("switched operation mid-line: %s -> %s ; line: %s\n", cur_entry_type, entry_type, line)
				return 0
			}
			cur_entry_type = entry_type
			if entry_type == MathEntryNumber {
				if len(problems) <= cur_idx {
					problems = append(problems, make([]uint, 0))
				}
				problems[cur_idx] = append(problems[cur_idx], entry_number)
			} else if entry_type == MathEntryOperation {
				val := uint(0)
				if entry_op == OperationTimes {
					val = performTimes(problems[cur_idx])
				} else if entry_op == OperationPlus {
					val = performPlus(problems[cur_idx])
				} else {
					fmt.Printf("invalid operation for entry %s\n", entry)
					return 0
				}
				fmt.Println(problems[cur_idx])
				fmt.Printf("problem %d %s -> %d\n", cur_idx, entry_op, val)
				password += val
			}
			cur_idx++
		}
	}
	return password
}

type Problem struct {
	entries   []uint
	operation OperationType
}

func isColEmpty(col string) bool {
	for idx := range col {
		if col[idx] != ' ' {
			return false
		}
	}
	return true
}
func parseProblemFromCols(yield func(Problem) bool) {
	// var prevCol *string = nil
	isStartOfProblem := true
	problem := Problem{make([]uint, 0), OperationNone}
	for col := range parseColumnsFromStdin {
		// prevCol = &col
		is_col_empty := isColEmpty(col)
		if isStartOfProblem {
			if is_col_empty {
				panic("empty col at start of problem")
			}
			op_candidate := col[len(col)-1]
			op := OperationNone
			switch op_candidate {
			case '*':
				op = OperationTimes
			case '+':
				op = OperationPlus
			default:
				panic(fmt.Sprintf("invalid operation %c in col %s", op_candidate, col))
			}
			problem.operation = op
			isStartOfProblem = false
		} else if is_col_empty {
			if problem.operation == OperationNone {
				panic(fmt.Sprintf("operation not set for problem %s\n", problem))
			}
			yield(problem)
			problem.entries = make([]uint, 0)
			problem.operation = OperationNone
			isStartOfProblem = true
			continue
		}
		candidateInt := strings.ReplaceAll(col[:len(col)-1], " ", "")
		val, err := strconv.Atoi(candidateInt)
		if err != nil {
			panic(err)
		}
		problem.entries = append(problem.entries, uint(val))
	}
	if !isStartOfProblem {
		yield(problem)
	}
}

func aoc6b() uint {
	// problems := make([][]uint, 0)
	password := uint(0)
	for problem := range parseProblemFromCols {
		fmt.Printf("--- %s\n", problem)
		result := uint(0)
		if problem.operation == OperationPlus {
			result = performPlus(problem.entries)
		} else if problem.operation == OperationTimes {
			result += performTimes(problem.entries)
		} else {
			fmt.Printf("invalid operation %s", problem.operation)
		}
		fmt.Printf("    -> %d\n", result)
		password += result
	}
	return password
}
