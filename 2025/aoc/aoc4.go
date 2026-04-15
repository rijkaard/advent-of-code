package main

import (
	"fmt"
)

type SlotType int

const (
	SlotEmpty SlotType = iota
	SlotRoll
	SlotRollAccessible
	// SlotUnused
)

var SlotToByte = map[SlotType]byte{
	SlotEmpty:          '.',
	SlotRoll:           '@',
	SlotRollAccessible: 'x',
	// SlotUnused: '@',
}

var ByteToSlot = makeByteToSlot()

func makeByteToSlot() map[byte]SlotType {
	out := make(map[byte]SlotType)
	for k, v := range SlotToByte {
		if _, hasKey := out[v]; hasKey {
			panic("key already present")
		}
		out[v] = k
	}
	return out
}

func boardToString(board [][]SlotType) (string, error) {
	out := ""
	boardHeight := len(board)
	if boardHeight == 0 {
		return "", nil
	}
	boardWidth := len(board[0])
	for row := range boardHeight {
		curLine := ""
		for col := range boardWidth {
			byteValue, hasKey := SlotToByte[board[row][col]]
			if !hasKey {
				return "", fmt.Errorf("unknown Slot %s\n:", board[row][col])
			}
			curLine += string(byteValue)
		}

		out += curLine + "\n"
	}
	return out, nil
}

func parseWarehouseLine(row string) ([]SlotType, error) {
	out := make([]SlotType, len(row))
	for pos := range row {
		slotValue, hasKey := ByteToSlot[row[pos]]
		if !hasKey {
			return out, fmt.Errorf("unknown %c", row[pos])
		}
		out[pos] = slotValue
	}
	return out, nil
}

func warehouseBoardFromStdin() ([][]SlotType, error) {
	board := make([][]SlotType, 0)
	for line := range parseLinesFromStdin {
		warehouse, err := parseWarehouseLine(line)
		if err != nil {
			return board, fmt.Errorf("couldn't parse line %s: %s", line, err)
		}
		board = append(board, warehouse)
	}
	if len(board) != 0 {
		width := len(board[0])
		for rowIdx := range len(board) - 1 {
			curLine := board[rowIdx+1]
			curLen := len(curLine)
			if width != curLen {
				return board, fmt.Errorf("invalid length for line %d: %d != %d", rowIdx+1, curLen, width)
			}
		}
	}
	return board, nil
}

var convolutionNeighbours = [][2]int{
	{-1, -1},
	{-1, +0},
	{-1, +1},
	{+0, -1},
	{+0, +1},
	{+1, -1},
	{+1, +0},
	{+1, +1},
}

func processBoard(board [][]SlotType, neighboursThreshold uint, removeAccessible bool) ([][]SlotType, uint, error) {
	var accessibleMarker SlotType = SlotRollAccessible
	if removeAccessible {
		accessibleMarker = SlotEmpty
	}
	accessibleCount := uint(0)
	boardHeight := len(board)
	if boardHeight == 0 {
		return nil, 0, nil
	}
	boardWidth := len(board[0])
	out := make([][]SlotType, boardHeight)
	for idx := range boardHeight {
		out[idx] = make([]SlotType, boardWidth)
	}
	for row := range boardHeight {
		for col := range boardWidth {
			boardValue := board[row][col]
			if boardValue == SlotEmpty {
				out[row][col] = SlotEmpty
				continue
			}
			if boardValue != SlotRoll {
				return nil, 0, fmt.Errorf("invalid slot type: %s", string(boardValue))
			}

			rollCount := uint(0)
			for convEntry := range len(convolutionNeighbours) {
				convRow, convCol := convolutionNeighbours[convEntry][0], convolutionNeighbours[convEntry][1]
				destRow, destCol := row+convRow, col+convCol
				if (destRow >= 0 && destRow < boardHeight) && (destCol >= 0 && destCol < boardWidth) && board[destRow][destCol] == SlotRoll {
					rollCount++
				}
			}
			if rollCount < neighboursThreshold {
				out[row][col] = accessibleMarker
				accessibleCount++
			} else {
				out[row][col] = SlotRoll
			}
		}
	}
	return out, accessibleCount, nil
}

func aoc4a(neighboursThreshold uint) uint {
	board, err := warehouseBoardFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(board)
	processedBoard, accessibleCount, err := processBoard(board, neighboursThreshold, false)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(boardToString(processedBoard))
	return accessibleCount
}

func aoc4b(neighboursThreshold uint) uint {
	password := uint(0)
	board, err := warehouseBoardFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(board)
	for {
		processedBoard, accessibleCount, err := processBoard(board, neighboursThreshold, true)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		fmt.Printf("accessible: %d\n", accessibleCount)
		if accessibleCount == 0 {
			break
		}
		password += accessibleCount
		board = processedBoard
	}
	return password
}
