package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TilePosition [2]int

func parsePair(line string) (TilePosition, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return TilePosition{0, 0}, fmt.Errorf("invalid line `%s`", line)
	}
	v1, err := strconv.Atoi(parts[0])
	if err != nil {
		return TilePosition{0, 0}, fmt.Errorf("invalid entry `%s` in line `%s`", parts[0], line)
	}
	v2, err := strconv.Atoi(parts[1])
	if err != nil {
		return TilePosition{0, 0}, fmt.Errorf("invalid entry `%s` in line `%s`", parts[1], line)
	}
	return TilePosition{v1, v2}, nil
}

func getPairsFromStdin() ([]TilePosition, error) {
	out := make([]TilePosition, 0)
	for line := range parseLinesFromStdin {
		if line == "" {
			break
		}
		pair, err := parsePair(line)
		if err != nil {
			return nil, err
		}
		out = append(out, pair)
	}
	return out, nil
}

func tileCoverage(p1, p2 TilePosition) int {
	d1 := p1[0] - p2[0]
	d1 = max(d1, -d1) + 1
	d2 := p1[1] - p2[1]
	d2 = max(d2, -d2) + 1
	return d1 * d2
}

func aoc9a() uint {
	pairs, err := getPairsFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(pairs)
	max_coverage := 0
	for fst_idx := range len(pairs) {
		for scn_idx := range len(pairs) - fst_idx - 1 {
			scn_idx += fst_idx + 1
			cur_coverage := tileCoverage(pairs[fst_idx], pairs[scn_idx])
			// fmt.Println(pairs[fst_idx], pairs[scn_idx], cur_coverage)
			max_coverage = max(max_coverage, cur_coverage)
		}
	}
	return uint(max_coverage)
}
