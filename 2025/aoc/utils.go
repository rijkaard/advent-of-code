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
