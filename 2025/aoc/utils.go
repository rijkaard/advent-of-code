package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func parseLinesFromStdin(yield func(string) bool) {
	reader := bufio.NewReader(os.Stdin)
	text := ""
	for true {
		curPart, isPrefix, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		text += string(curPart)
		if isPrefix {
			continue
		}
		if !yield(text) {
			break
		}
		text = ""
		if err == io.EOF {
			break
		}
	}
}

func parseColumnsFromStdin(yield func(string) bool) {
	lines := make([]string, 0)
	for line := range parseLinesFromStdin {
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		panic("no lines")
	}
	expected_length := len(lines[0])
	for _, line := range lines[1:] {
		if len(line) != expected_length {
			panic("inconsistent length")
		}
	}
	curcol := make([]byte, len(lines))
	for idx := range expected_length {
		for pos, line := range lines {
			curcol[pos] = line[idx]
		}
		if !yield(string(curcol)) {
			break
		}
	}
}
