package main

import (
	"bufio"
	"os"
)

func parseLinesFromStdin(yield func(string) bool) {
	reader := bufio.NewReader(os.Stdin)
	text := ""
	for true {
		curPart, isPrefix, _ := reader.ReadLine()
		text += string(curPart)
		if isPrefix {
			continue
		}
		if len(text) == 0 {
			break
		}
		if !yield(text) {
			break
		}
		text = ""
	}
}
