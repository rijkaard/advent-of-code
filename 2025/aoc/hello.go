// hello.go

// package aoc25
package main

import (
	"fmt"
	"os"
)

// type FunctionMap map[string]func(...any)
// var x FunctionMap = FunctionMap{
// 		"aoc1": func(...any){ return aoc1(50, 100)}
// }

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command>", os.Args[0])
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "aoc1a":
		fmt.Println(aoc1a(50, 100))
	case "aoc1b":
		fmt.Println(aoc1b(50, 100))
	case "aoc2a":
		fmt.Println(aoc2a())
	case "aoc2b":
		fmt.Println(aoc2b())
	default:
		fmt.Printf("Invalid command: %s", cmd)
	}
}
