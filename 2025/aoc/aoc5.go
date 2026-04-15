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

func intersectRanges(first, second IDRange) (IDRange, bool) {
	if first[0] > second[0] {
		first, second = second, first
	}
	if first[1] < second[0] {
		return IDRange{0, 0}, false
	}
	return IDRange{first[0], max(first[1], second[1])}, true
}

func mergeRange(new IDRange, ranges []IDRange) []IDRange {
	out := make([]IDRange, len(ranges), len(ranges)+1)
	copy(out, ranges)
	ranges = out
	for {
		mergerHappened := false
		for idx := range len(ranges) {
			intersection, hasIntersected := intersectRanges(new, ranges[idx])
			if hasIntersected {
				// fmt.Printf("intersected: %s + %s = %s", new, ranges[idx], intersection)
				mergerHappened = true
				ranges = append(ranges[:idx], ranges[idx+1:]...)
				new = intersection
				break
			}
		}
		if !mergerHappened {
			ranges = append(ranges, new)
			break
		}
	}
	return ranges
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

func aoc5b() uint {
	password := uint(0)

	ranges := make([]IDRange, 0)

	for line := range parseLinesFromStdin {
		if line == "" {
			break
		}
		start, end, err := parseRange(line)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		ranges = mergeRange(IDRange{start, end}, ranges)
	}
	fmt.Println(ranges)
	for rngIdx := range ranges {
		rng := ranges[rngIdx]
		password += rng[1] - rng[0] + 1
	}
	return password
}
