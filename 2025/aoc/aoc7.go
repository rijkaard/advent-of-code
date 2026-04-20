package main

import (
	"fmt"
	"strings"
)

type BeamSplitterEntryType uint

const (
	BeamSplitterSource BeamSplitterEntryType = iota
	BeamSplitterEmpty
	BeamSplitterSplitter
	BeamSplitterBeam
)

func maybePrintBeamSplitterRow(row []BeamSplitterEntryType) string {
	val, err := printBeamSplitterRow(row)
	if err != nil {
		return fmt.Sprint(row)
	}
	return val
}

func printBeamSplitterRow(row []BeamSplitterEntryType) (string, error) {
	out := ""
	const accepted = ".^S|"
	for _, state := range row {
		if strings.Contains(accepted, string(state)) {
			out = out + "."
		} else {
		}
		switch state {
		case BeamSplitterEmpty:
			out += "."
		case BeamSplitterSource:
			out += "S"
		case BeamSplitterSplitter:
			out += "^"
		case BeamSplitterBeam:
			out += "|"
		default:
			return "", fmt.Errorf("Invalid state %d in row `%s`", state, row)
		}
	}
	return out, nil
}

func parseBeamSplitterRow(row string) ([]BeamSplitterEntryType, error) {
	out := make([]BeamSplitterEntryType, 0)
	for _, chr := range row {
		switch chr {
		case '.':
			out = append(out, BeamSplitterEmpty)
		case '^':
			out = append(out, BeamSplitterSplitter)
		case 'S':
			out = append(out, BeamSplitterSource)
		default:
			return nil, fmt.Errorf("invalid character `%c` in row `%s`", chr, row)
		}
	}
	return out, nil
}

func parseBeamSplitterRows(yield func([]BeamSplitterEntryType) bool) {
	for line := range parseLinesFromStdin {
		row, err := parseBeamSplitterRow(line)
		if err != nil {
			panic(fmt.Sprintf("invalid row `%s`", line))
		}
		if !yield(row) {
			break
		}
	}
}

func ensureSingleSource(state []BeamSplitterEntryType) error {
	n_sources := uint(0)
	for _, chr := range state {
		if chr == BeamSplitterSource {
			n_sources++
		}
	}
	if n_sources != 1 {
		return fmt.Errorf("found %d sources, expected 1 in row %s", n_sources, maybePrintBeamSplitterRow(state))
	}
	return nil
}

func beamSplitterStep(state, next []BeamSplitterEntryType) (uint, []BeamSplitterEntryType, error) {
	split_count := uint(0)
	if len(state) != len(next) {
		panic("incompatible lengths")
	}
	out := make([]BeamSplitterEntryType, len(state))
	for pos := range state {
		out[pos] = BeamSplitterEmpty
	}
	for pos, cur_state := range state {
		cur_next := next[pos]
		if cur_state == BeamSplitterSource || cur_state == BeamSplitterBeam {
			if cur_next == BeamSplitterEmpty {
				out[pos] = BeamSplitterBeam
			} else if cur_next == BeamSplitterSplitter {
				split_count++
				if pos > 0 {
					if next[pos-1] != BeamSplitterEmpty {
						return 0, nil, fmt.Errorf("unexpected non-empty state %d in next: %s", next[pos-1], maybePrintBeamSplitterRow(next))
					}
					out[pos-1] = BeamSplitterBeam
				}
				if pos < len(state)-1 {
					if next[pos+1] != BeamSplitterEmpty {
						return 0, nil, fmt.Errorf("unexpected non-empty state %d in next: %s", next[pos+1], maybePrintBeamSplitterRow(next))
					}
					out[pos+1] = BeamSplitterBeam
				}
			} else {
				return 0, nil, fmt.Errorf("unexpected combination - state: %d next: %d", cur_state, cur_next)
			}
		} else if !(cur_state == BeamSplitterSplitter || cur_state == BeamSplitterEmpty) {
			return 0, nil, fmt.Errorf("unexpected state %d in state`%s`", cur_state, maybePrintBeamSplitterRow(state))
		}
	}
	return split_count, out, nil
}

func aoc7a() uint {
	password := uint(0)
	rows := make([][]BeamSplitterEntryType, 0)
	var cur []BeamSplitterEntryType = nil
	for cur_row := range parseBeamSplitterRows {
		if len(cur_row) == 0 {
			break
		}
		rows = append(rows, cur_row)
		// fmt.Printf("nxt %s\n", maybePrintBeamSplitterRow(cur_row))
		if cur == nil {
			err := ensureSingleSource(cur_row)
			if err != nil {
				fmt.Println(err)
				return 0
			}
			cur = cur_row
		} else {
			if len(cur) != len(cur_row) {
				fmt.Printf("invalid length for current row `%s`", maybePrintBeamSplitterRow(cur_row))
			}
			split_count, candidate, err := beamSplitterStep(cur, cur_row)
			if err != nil {
				fmt.Println(err)
				return 0
			}
			// fmt.Println(split_count)
			password += split_count
			cur = candidate
		}
		// fmt.Printf("--> %s\n", maybePrintBeamSplitterRow(cur))
	}
	return password
}
