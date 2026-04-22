package main

import (
	"fmt"
	"strconv"
	"strings"
)

type BoardRanges struct {
	perRow [][]Range[uint]
	perCol [][]Range[uint]
}

func initBoardRanges(n uint) BoardRanges {
	brdrng := BoardRanges{
		make([][]Range[uint], n),
		make([][]Range[uint], n),
	}
	for idx := range n {
		brdrng.perCol[idx] = make([]Range[uint], 0)
		brdrng.perRow[idx] = make([]Range[uint], 0)
	}
	return brdrng
}

type TilePosition [2]int

func parsePair(line string) (TilePosition, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return TilePosition{0, 0}, fmt.Errorf("invalid line `%s`", line)
	}
	v1, err := strconv.Atoi(parts[0])
	if err != nil {
		return TilePosition{0, 0}, fmt.Errorf("invalid entry `%s` in line `%s`", parts[0], line)
	}
	v2, err := strconv.Atoi(parts[1])
	if err != nil {
		return TilePosition{0, 0}, fmt.Errorf("invalid entry `%s` in line `%s`", parts[1], line)
	}
	return TilePosition{v1, v2}, nil
}

func getPairsFromStdin() ([]TilePosition, error) {
	out := make([]TilePosition, 0)
	for line := range parseLinesFromStdin {
		if line == "" {
			break
		}
		pair, err := parsePair(line)
		if err != nil {
			return nil, err
		}
		out = append(out, pair)
	}
	return out, nil
}

func tileCoverage(p1, p2 TilePosition) int {
	d1 := p1[0] - p2[0]
	d1 = max(d1, -d1) + 1
	d2 := p1[1] - p2[1]
	d2 = max(d2, -d2) + 1
	return d1 * d2
}

func aoc9a() uint {
	pairs, err := getPairsFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(pairs)
	max_coverage := 0
	for fst_idx := range len(pairs) {
		for scn_idx := range len(pairs) - fst_idx - 1 {
			scn_idx += fst_idx + 1
			cur_coverage := tileCoverage(pairs[fst_idx], pairs[scn_idx])
			max_coverage = max(max_coverage, cur_coverage)
		}
	}
	return uint(max_coverage)
}

func printBoard(board [][]bool) string {
	all_strings := make([]string, len(board))
	for _, row := range board {
		s := ""
		for _, val := range row {
			if val {
				s += "#"
			} else {
				s += "."
			}
		}
		all_strings = append(all_strings, s)
	}
	return strings.Join(all_strings, "\n")
}

type Board struct {
	board   [][]bool
	max_row int
	max_col int
}

func fillBoard(board [][]bool) {
	for _, row := range board {
		col := 0
		for col < len(row) {
			val := row[col]
			if val && col+1 < len(row) && !row[col+1] {
				for idx := col + 1; idx < len(row); idx++ {
					if row[idx] {
						col = idx
						break
					}
					row[idx] = true
				}
			} else {
				col++
			}
		}
	}
}

func makeContourBoard(pairs []TilePosition) (Board, error) {
	max_col, max_row := 0, 0
	for _, pair := range pairs {
		max_col = max(max_col, pair[0])
		max_row = max(max_row, pair[1])
	}

	board := make([][]bool, max_row+1)
	for idx := range len(board) {
		board[idx] = make([]bool, max_col+1)
	}
	prev_pair := pairs[0]
	for _, pair := range append(pairs, pairs[0]) {
		col, row := pair[0], pair[1]
		board[row][col] = true
		if prev_pair[0] == col {
			start := min(prev_pair[1], row)
			end := max(prev_pair[1], row)
			for idx := start; idx <= end; idx++ {
				board[idx][col] = true
			}
		} else if prev_pair[1] == row {
			start := min(prev_pair[0], col)
			end := max(prev_pair[0], col)
			for idx := start; idx <= end; idx++ {
				board[row][idx] = true
			}
		} else {
			return Board{nil, 0, 0}, fmt.Errorf("subsequent pairs not aligned: %v %v\n", prev_pair, pair)
		}
		prev_pair = pair
	}
	return Board{board, max_row, max_col}, nil
}

func makeDottedBoard(pairs []TilePosition) (Board, error) {
	max_col, max_row := 0, 0
	for _, pair := range pairs {
		max_col = max(max_col, pair[0])
		max_row = max(max_row, pair[1])
	}

	board := make([][]bool, max_row+1)
	for idx := range len(board) {
		board[idx] = make([]bool, max_col+1)
	}
	prev_pair := pairs[0]
	for _, pair := range append(pairs, pairs[0]) {
		col, row := pair[0], pair[1]
		board[row][col] = true
		if prev_pair[0] == col {
			start := min(prev_pair[1], row)
			end := max(prev_pair[1], row)
			for idx := start; idx <= end; idx++ {
				board[idx][col] = true
			}
		} else if prev_pair[1] == row {
			start := min(prev_pair[0], col)
			end := max(prev_pair[0], col)
			board[row][start] = true
			board[row][end] = true
		} else {
			return Board{nil, 0, 0}, fmt.Errorf("subsequent pairs not aligned: %v %v\n", prev_pair, pair)
		}
		prev_pair = pair
	}
	return Board{board, max_row, max_col}, nil
}

func makeFilledBoard(pairs []TilePosition) (Board, error) {
	board, err := makeContourBoard(pairs)
	if err != nil {
		return board, err
	}
	fmt.Println(printBoard(board.board))
	fillBoard(board.board)
	return board, nil
}

func makeBoardAndRangesv0(pairs []TilePosition) (Board, BoardRanges, error) {
	board, err := makeContourBoard(pairs)
	if err != nil {
		return board, BoardRanges{nil, nil}, err
	}
	fillBoard(board.board)
	fmt.Println(printBoard(board.board))
	boardRanges := initBoardRanges(uint(int(max(board.max_row, board.max_col) + 1)))
	for row_idx, row := range board.board {
		fmt.Println(printBoard(board.board[row_idx : row_idx+1]))
		col := 0
		range_start := -1
		for col < len(row) {
			val := row[col]
			if val && range_start < 0 {
				range_start = col
			} else if !val && range_start >= 0 {
				range_end := col - 1
				boardRanges.perRow[row_idx] = mergeRange(Range[uint]{uint(range_start), uint(range_end)}, boardRanges.perRow[row_idx], true)
				fmt.Println("merged", row_idx, range_start, range_end, boardRanges.perRow[row_idx])
				range_start = -1
			}
			col++
		}
		if range_start >= 0 {
			range_end := col - 1
			boardRanges.perRow[row_idx] = mergeRange(Range[uint]{uint(range_start), uint(range_end)}, boardRanges.perRow[row_idx], true)
			fmt.Println("merged", row_idx, range_start, range_end, boardRanges.perRow[row_idx])
		}
	}
	return board, boardRanges, nil
}

func makeBoardAndRangesv1(pairs []TilePosition) (Board, BoardRanges, error) {
	board, err := makeContourBoard(pairs)
	if err != nil {
		return board, BoardRanges{nil, nil}, err
	}
	fillBoard(board.board)
	fmt.Println(printBoard(board.board))
	boardRanges := initBoardRanges(uint(int(max(board.max_row, board.max_col) + 1)))
	for row_idx, row := range board.board {
		fmt.Println(printBoard(board.board[row_idx : row_idx+1]))
		col := 0
		range_start := -1
		for col < len(row) {
			val := row[col]
			if val && range_start < 0 {
				range_start = col
			} else if !val && range_start >= 0 {
				range_end := col - 1
				boardRanges.perRow[row_idx] = mergeRange(Range[uint]{uint(range_start), uint(range_end)}, boardRanges.perRow[row_idx], true)
				fmt.Println("merged", row_idx, range_start, range_end, boardRanges.perRow[row_idx])
				range_start = -1
			}
			col++
		}
		if range_start >= 0 {
			range_end := col - 1
			boardRanges.perRow[row_idx] = mergeRange(Range[uint]{uint(range_start), uint(range_end)}, boardRanges.perRow[row_idx], true)
			fmt.Println("merged", row_idx, range_start, range_end, boardRanges.perRow[row_idx])
		}
	}
	return board, boardRanges, nil
}

func makeBoardAndRanges(pairs []TilePosition) (Board, BoardRanges, error) {
	board, err := makeDottedBoard(pairs)
	if err != nil {
		return board, BoardRanges{nil, nil}, err
	}
	boardRanges := initBoardRanges(uint(int(max(board.max_row, board.max_col) + 1)))
	for row_idx, row := range board.board {
		col := 0
		range_start := -1
		range_end := -1
		for col < len(row) {
			val := row[col]
			if val && range_start < 0 {
				range_start = col
			} else if val && range_start >= 0 {
				range_end = col
			}
			col++
		}
		if range_start >= 0 {
			boardRanges.perRow[row_idx] = mergeRange(Range[uint]{uint(range_start), uint(range_end)}, boardRanges.perRow[row_idx], true)
			range_start = -1
			range_end = -1
		} else if range_end >= 0 {
			panic(fmt.Sprintf("unexpected end range in line %s", printBoard(board.board[row_idx:row_idx+1])))
		}
	}
	return board, boardRanges, nil
}

func aoc9bv0() uint {
	pairs, err := getPairsFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if len(pairs) == 0 {
		fmt.Println("empty list of pairs")
		return 0
	}
	board, err := makeFilledBoard(pairs)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println(printBoard(board.board))
	fmt.Println(pairs)
	return 0
}

func isFullyContained(p1, p2 TilePosition, board_ranges BoardRanges) bool {
	min_row := min(p1[1], p2[1])
	max_row := max(p1[1], p2[1])

	min_col := min(p1[0], p2[0])
	max_col := max(p1[0], p2[0])
	col_range := Range[uint]{uint(min_col), uint(max_col)}
	for row_idx := min_row; row_idx <= max_row; row_idx++ {
		if !rangeInRanges(col_range, board_ranges.perRow[row_idx]) {
			return false
		}
	}
	return true
}

func aoc9b() uint {
	pairs, err := getPairsFromStdin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if len(pairs) == 0 {
		fmt.Println("empty list of pairs")
		return 0
	}
	_, ranges, err := makeBoardAndRanges(pairs)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	password := uint(0)
	for fst_idx, fst := range pairs {
		for scn_idx := range len(pairs) - fst_idx - 1 {
			scn_idx += fst_idx + 1
			scn := pairs[scn_idx]
			if isFullyContained(fst, scn, ranges) {
				// fmt.Println("fully contained", fst, scn)
				password = max(password, uint(tileCoverage(fst, scn)))
			}
		}
	}
	return password
}
