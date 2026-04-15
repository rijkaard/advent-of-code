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
			// fmt.Printf("line too large")
			// return
		}
		if len(text) == 0 {
			break
		}
		// instr, err := parseInstruction(string(text))
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		//fmt.Println(instr)
		if !yield(text) {
			break
		}
		text = ""
	}
}
