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
		fmt.Printf("--- %s\n", line)
		cur_entry_type := MathEntryNone
		cur_idx := 0
		for entry := range yieldEntriesFromLine(line) {
			// fmt.Printf("    %s\n", entry)
			entry_type, entry_number, entry_op, err := parseEntry(entry)
			// fmt.Printf("    -> %s %d %s %s\n", entry_type, entry_number, entry_op, err)
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
