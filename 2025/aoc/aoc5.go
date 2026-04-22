package main

import (
	"fmt"
	"strconv"
)

func aoc5a() uint {
	password := uint(0)
	state := StateWantsRanges

	ranges := make([]Range[uint], 0)

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
			ranges = append(ranges, Range[uint]{start, end})
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
			if matchesRanges[uint](ingredient, ranges) {
				fmt.Printf("matches: %d\n", ingredient)
				password++
			}
		} else if state == StateWantsEOF {
			fmt.Printf("expected EOF, not %s\n", line)
		}
	}
	return password
}

func aoc5b() uint {
	password := uint(0)

	ranges := make([]Range[uint], 0)

	for line := range parseLinesFromStdin {
		if line == "" {
			break
		}
		start, end, err := parseRange(line)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		ranges = mergeRange(Range[uint]{start, end}, ranges, false)
	}
	fmt.Println(ranges)
	for rngIdx := range ranges {
		rng := ranges[rngIdx]
		password += rng[1] - rng[0] + 1
	}
	return password
}
