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

func getMaxDigitByte(text string) (byte, error) {
	res, err := getMaxDigit(text)
	if err != nil {
		return 0, err
	}
	return '0' + byte(res), nil
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

func compareDigits(fst []byte, scn []byte) (int, error) {
	if len(fst) != len(scn) {
		return 0, fmt.Errorf("arrays have different lengths")
	}
	for idx := range len(fst) {
		comp := int16(fst[idx]) - int16(scn[idx])
		if comp != 0 {
			if comp < 0 {
				return -1, nil
			}
			return 1, nil
		}
	}
	return 0, nil
}

func setDigitsToValue(dst []byte, value byte) {
	for idx := range len(dst) {
		dst[idx] = value
	}
}

func setDigits(src []byte, dst []byte) {
	if len(src) != len(dst) {
		panic("different size")
	}
	for idx := range len(src) {
		dst[idx] = src[idx]
	}
}

func innerGetBankMaxJoltageN(digits []byte, curMax []byte, bank string, n uint) {
	lenDigits := uint(len(digits))
	lenBank := uint(len(bank))

	if lenBank < lenDigits-n {
		return
	}
	if n == lenDigits {
		compareResult, err := compareDigits(curMax, digits)
		if err != nil {
			panic("should not happen here")
		}
		if compareResult < 0 {
			// fmt.Printf("%s -> %s", curMax, digits)
			setDigits(digits, curMax)
		}
		return
	}
	remaining := lenDigits - n - 1
	// for idx := uint(0); idx < lenBank; idx++ {
	for idx := uint(0); idx < lenBank-remaining; idx++ {
		digits[n] = bank[idx]
		innerGetBankMaxJoltageN(digits, curMax, bank[idx+1:], n+1)
	}
}

func innerGetBankMaxJoltageNv3(digitsPtr *[]byte, curMaxPtr *[]byte, bankPtr *string, bankBase uint, n uint) {
	digits := *digitsPtr
	curMax := *curMaxPtr
	bank := *bankPtr
	lenDigits := uint(len(digits))
	lenBank := uint(len(bank))

	if lenBank < lenDigits-n {
		return
	}
	if n == lenDigits {
		compareResult, err := compareDigits(curMax, digits)
		if err != nil {
			panic("should not happen here")
		}
		if compareResult < 0 {
			// fmt.Printf("%s -> %s", curMax, digits)
			setDigits(digits, curMax)
		}
		return
	}
	remaining := lenDigits - n - 1
	// for idx := uint(0); idx < lenBank; idx++ {
	for idx := bankBase; idx < lenBank-remaining; idx++ {
		digits[n] = bank[idx]
		// if digits[n] >= curMax[n] {
		innerGetBankMaxJoltageNv3(digitsPtr, curMaxPtr, bankPtr, idx+1, n+1)
		// }
	}
}

func innerGetBankMaxJoltageNv4(digitsPtr *[]byte, curMaxPtr *[]byte, bankPtr *string, bankEnd uint, n uint) {
	digits := *digitsPtr
	curMax := *curMaxPtr
	bank := *bankPtr
	lenDigits := uint(len(digits))
	lenBank := uint(len(bank))

	if lenBank < lenDigits-n {
		return
	}
	if n == 0 {
		maxDigit, err := getMaxDigitByte(bank[:bankEnd])
		if err != nil {
			panic(fmt.Sprintf("couldn't get max digit: %s", err))
		}
		digits[0] = maxDigit
		compareResult, err := compareDigits(curMax, digits)
		if err != nil {
			panic("should not happen here")
		}
		if compareResult < 0 {
			// fmt.Printf("%s -> %s", curMax, digits)
			setDigits(digits, curMax)
		}
		return
	}
	for idx := bankEnd - 1; idx >= n; idx-- {
		digits[n] = bank[bankEnd-idx-1]
		innerGetBankMaxJoltageNv4(digitsPtr, curMaxPtr, bankPtr, idx-1, n+1)
	}
}

func innerGetBankMaxJoltageNv5(curMax []byte, bank string, n uint) {
	lenDigits := len(curMax)
	lenBank := len(bank)

	if int(n) == lenDigits {
		return
	}
	remaining := lenDigits - int(n) - 1
	maxIdx := 0
	max := byte(0)
	for idx := 0; idx < lenBank-remaining; idx++ {
		if bank[idx] > max {
			max = bank[idx]
			maxIdx = idx
		}
	}
	curMax[n] = max
	// fmt.Printf("setting %d to %c (%d)", n, max, maxIdx)
	innerGetBankMaxJoltageNv5(curMax, bank[maxIdx+1:], n+1)
}

func innerGetBankMaxJoltageNv6(curMax []byte, bank string) {
	lenDigits := len(curMax)
	lenBank := len(bank)
	bankStart := 0

	for n := range lenDigits {
		remaining := lenDigits - n - 1
		maxIdx := 0
		max := byte(0)
		for idx := 0; bankStart+idx < lenBank-remaining; idx++ {
			actualIdx := bankStart + idx
			if bank[actualIdx] > max {
				max = bank[actualIdx]
				maxIdx = idx
			}
		}
		curMax[n] = max
		bankStart += maxIdx + 1
	}
}

func bootstrapMax(bank string, n uint) ([]byte, error) {
	bankLen := len(bank)
	out := make([]byte, n, n)
	if bankLen < int(n) {
		return out, fmt.Errorf("invalid error")
	}
	currMax, err := getMaxDigitByte(bank[:bankLen-int(n)-1])
	if err != nil {
		return out, err
	}
	out[0] = currMax
	for idx := range int(n) - 1 {
		out[idx+1] = bank[bankLen-int(n)+idx+1]
	}
	return out, nil
}

func getBankMaxJoltageN(bank string, n uint) (string, error) {
	// instructions:
	// - use array of characters (or string, but we need to mutate them)
	// - implement a comparison between strings; probably strings.Compare is enough
	// - call getBankMaxJoltageN recursively
	max := make([]byte, n, n)
	setDigitsToValue(max, 0)
	// innerGetBankMaxJoltageNv5(max, bank, 0)
	innerGetBankMaxJoltageNv6(max, bank)
	return string(max), nil
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
	}
	return password
}

func aoc3b(n uint) uint {
	var password uint = 0
	for bank := range parseBatteriesFromStdin {
		maxJoltage, err := getBankMaxJoltageN(bank, n)
		if err != nil {
			fmt.Printf("failed processing bank %s: %s", bank, err)
		}
		// fmt.Printf("bank: %s -> maxJoltage: %s\n", bank, maxJoltage)
		maxJoltageNumeric, err := strconv.Atoi(string(maxJoltage))
		if err != nil {
			fmt.Printf("invalid joltage: %s", maxJoltage)
			return 0
		}
		password += uint(maxJoltageNumeric)
	}
	return password
}
