package main

import (
	"fmt"
	"strconv"
)

type IDRange [2]uint
type InputStateType int

const (
	StateWantsRanges InputStateType = iota
	StateWantsIngredients
	StateWantsEOF
)

func matchesRange(n uint, rng IDRange) bool {
	return n >= rng[0] && n <= rng[1]
}

func matchesRanges(n uint, ranges []IDRange) bool {
	for _, rng := range ranges {
		if matchesRange(n, rng) {
			return true
		}
	}
	return false
}

func aoc5a() uint {
	password := uint(0)
	state := StateWantsRanges

	ranges := make([]IDRange, 0)

	for line := range parseLinesFromStdin {
		if state == StateWantsRanges {
			if line == "" {
				state = StateWantsIngredients
				continue
			}
			start, end, err := parseRange(line)
			if err != nil {
				fmt.Println(err)
				return 0
			}
			ranges = append(ranges, IDRange{start, end})
		} else if state == StateWantsIngredients {
			if line == "" {
				state = StateWantsEOF
				continue
			}
			ingrValue, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println(err)
				return 0
			}
			ingredient := uint(ingrValue)
			if matchesRanges(ingredient, ranges) {
				fmt.Printf("matches: %d\n", ingredient)
				password++
			}
		} else if state == StateWantsEOF {
			fmt.Printf("expected EOF, not %s\n", line)
		}
	}
	return password
}
