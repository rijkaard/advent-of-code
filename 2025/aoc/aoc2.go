package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func isValidID(candidate string) bool {
	// invalid if starts with 0
	if len(candidate) == 0 || candidate[0] == '0' {
		return false
	}
	// invalid it's made of a subsequence repeated twice
	if len(candidate)%2 == 1 {
		// odd digits cannot be repeated
		return true
	}
	halfway := len(candidate) / 2
	idx := 0
	for idx < halfway {
		if candidate[idx] != candidate[halfway+idx] {
			return true
		}
		idx++
	}
	return false
}

func extractInvalidIDInRange(start uint, end uint) <-chan uint {
	out := make(chan uint)
	go func() {
		current := start
		for current <= end {
			if !isValidID(strconv.Itoa(int(current))) {
				// fmt.Println(current)
				out <- current
			}
			current++
		}
		close(out)
	}()
	return out
}

func parseRange(candidate string) (start uint, end uint, err error) {
	parts := strings.Split(candidate, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Invalid range: %s", candidate)
	}
	fst, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	scn, err := strconv.Atoi(parts[1])
	if err != nil {
		return uint(fst), 0, err
	}
	return uint(fst), uint(scn), nil
}

func parseRangesFromStdin(yield func(uint, uint) bool) {
	reader := bufio.NewReader(os.Stdin)
	for true {
		text, err := reader.ReadString(',')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\n", "")
		if len(text) == 0 {
			break
		}
		if text[len(text)-1] == ',' {
			text = text[:len(text)-1]
		}
		start, end, err := parseRange(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		if start > end {
			fmt.Printf("invalid start < end: %s\n", text)
			return
		}
		if !yield(start, end) {
			break
		}
	}
}

func aoc2a() uint {
	var password uint = 0
	for start, end := range parseRangesFromStdin {
		// fmt.Printf("%d -> %d\n", start, end)
		for invalid := range extractInvalidIDInRange(start, end) {
			// fmt.Printf("invalid: %d\n", invalid)
			password += invalid
		}
	}
	return password
}
