package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var sentinel = "asd"

// type Instruction struct {
// }
type Instruction int

func turn(start int, inst Instruction, nTicks int) int {
	return (start + int(inst)) % nTicks
}

// func turnLeft(start int, n int, nTicks int) int {
// 	return (start - n) % nTicks
// }

// func turnRight(start int, n int, nTicks int) int {
// 	return (start + n) % nTicks
// }

func parseInstruction(line string) (Instruction, error) {
	if len(line) == 0 {
		return 0, fmt.Errorf("Empty line")
	}
	z := strings.ToLower(line[:1])[0]
	var out Instruction
	switch z {
	case 'l':
		out = -1
	case 'r':
		out = 1
	default:
		return 0, fmt.Errorf("Invalid command: %s", line)
	}
	ticks, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, fmt.Errorf("Invalid command: %s", line)
	}
	return out * Instruction(ticks), nil
}

func parseFromStdin(yield func(Instruction) bool) {
	reader := bufio.NewReader(os.Stdin)
	for true {
		text, isPrefix, _ := reader.ReadLine()
		if isPrefix {
			fmt.Printf("line too large")
			return
		}
		if len(text) == 0 {
			break
		}
		instr, err := parseInstruction(string(text))
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(instr)
		if !yield(instr) {
			break
		}
	}
}

func aoc1a(start int, nTicks int) int {
	password := 0
	current := start
	for instr := range parseFromStdin {
		current = turn(current, instr, nTicks)
		// fmt.Println(current)
		if current == 0 {
			password++
		}
	}
	return password
}
