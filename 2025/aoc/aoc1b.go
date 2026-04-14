package main

func turnRightWithZero(start int, turnTicks uint, nTicks int) (end int, nCrosses uint) {
	if turnTicks == 0 {
		panic("nope")
	}
	nCrosses = turnTicks / uint(nTicks)
	turnTicks -= nCrosses * uint(nTicks)
	end = turn(start, Instruction(int(turnTicks)), nTicks)
	if end < start {
		nCrosses++
	}
	return end, nCrosses
}

func turnLeftWithZero(start int, turnTicks uint, nTicks int) (end int, nCrosses uint) {
	revEnd, nCrosses := turnRightWithZero((nTicks-start)%nTicks, turnTicks, nTicks)
	return (nTicks - revEnd) % nTicks, nCrosses
}

func turnWithZero(start int, inst Instruction, nTicks int) (int, uint) {
	if inst < 0 {
		return turnLeftWithZero(start, uint(-int(inst)), nTicks)
	}
	return turnRightWithZero(start, uint(inst), nTicks)
}

func aoc1b(start int, nTicks int) uint {
	var nCrosses uint
	var password uint = 0
	current := start
	for instr := range parseFromStdin {
		current, nCrosses = turnWithZero(current, instr, nTicks)
		// fmt.Printf("(%d) -> %d %b\n", instr, current, nCrosses)
		password += nCrosses
	}
	return password
}
