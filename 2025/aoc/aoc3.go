package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func parseBatteryBank(text string) (string, error) {
	textLen := len(text)
	if textLen == 0 {
		return "", fmt.Errorf("invalid zero-sized bank")
	}
	for start := 0; start < textLen; start++ {
		if !isDigit(text[start]) {
			return "", fmt.Errorf("invalid digit %c in bank: %s", text[start], text)
		}
	}
	return text, nil
}

func digitToUInt(c byte) (uint, error) {
	if c < '0' || c > '9' {
		return 0, fmt.Errorf("invalid digit: %c", c)
	}
	return uint(c - '0'), nil
}

func getMaxDigit(text string) (uint, error) {
	if len(text) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	// fmt.Println(text)
	var out uint = 0
	for pos, _ := range text {
		digit, err := digitToUInt(text[pos])
		if err != nil {
			return 0, err
		}
		out = max(out, digit)
	}
	// fmt.Printf("rune %d byte %c\n", out, byte(out))
	return out, nil
}

func getBankMaxJoltage(text string) (uint, error) {
	// 5156910691 -> 99 largest number possible with 2 digits in the order they appear
	textLen := len(text)
	var curMax uint = 0
	var candidateMax uint
	for start := 0; start < textLen-1; start++ {
		head, err := strconv.Atoi(text[start : start+1])
		if err != nil {
			return 0, err
		}
		tail, err := getMaxDigit(text[start+1:])
		if err != nil {
			return 0, err
		}
		// fmt.Printf("head: %d tail: %d\n", head, tail)
		candidateMax = 10*uint(head) + uint(tail)
		curMax = max(curMax, candidateMax)
	}
	return curMax, nil
}

func parseBatteriesFromStdin(yield func(string) bool) {
	reader := bufio.NewReader(os.Stdin)
	for true {
		text, isPrefix, _ := reader.ReadLine()
		if isPrefix {
			fmt.Printf("line too large")
			return
		}
		if len(text) == 0 {
			break
		}
		bank, err := parseBatteryBank(string(text))
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(instr)
		if !yield(bank) {
			break
		}
	}
}

func aoc3a() uint {
	var password uint = 0
	for bank := range parseBatteriesFromStdin {
		maxJoltage, err := getBankMaxJoltage(bank)
		if err != nil {
			fmt.Printf("failed processing bank %s: %s", bank, err)
		}
		fmt.Printf("bank: %s -> maxJoltage: %d\n", bank, maxJoltage)
		password += maxJoltage
		// for invalid := range extractInvalidIDInRange(start, end, isValidIDV2) {
		// 	fmt.Printf("invalid: %d\n", invalid)
		// 	password += invalid
		// }
	}
	return password
}
