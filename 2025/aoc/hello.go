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
	case "aoc3a":
		fmt.Println(aoc3a())
	case "aoc3b":
		fmt.Println(aoc3b(12))
	case "aoc4a":
		fmt.Println(aoc4a(4))
	case "aoc4b":
		fmt.Println(aoc4b(4))
	case "aoc5a":
		fmt.Println(aoc5a())
	case "aoc5b":
		fmt.Println(aoc5b())
	case "aoc6a":
		fmt.Println(aoc6a())
	case "aoc6b":
		fmt.Println(aoc6b())
	default:
		fmt.Printf("Invalid command: %s", cmd)
	}
}
