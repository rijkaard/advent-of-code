package main

type UInteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}
type Integer interface {
	int8 | int16 | int32 | int64
}
type Float interface {
	float32 | float64
}

type RangeEntryType interface {
	UInteger | Integer | Float
}
type Range[T RangeEntryType] [2]T
type InputStateType int

const (
	StateWantsRanges InputStateType = iota
	StateWantsIngredients
	StateWantsEOF
)

func matchesRange[T RangeEntryType](n T, rng Range[T]) bool {
	return n >= rng[0] && n <= rng[1]
}

func matchesRanges[T RangeEntryType](n T, ranges []Range[T]) bool {
	for _, rng := range ranges {
		if matchesRange(n, rng) {
			return true
		}
	}
	return false
}

func intersectRanges[T RangeEntryType, R Range[T]](first, second R) (R, bool) {
	if first[0] > second[0] {
		first, second = second, first
	}
	if first[1] < second[0] {
		return R{0, 0}, false
	}
	return R{first[0], max(first[1], second[1])}, true
}

func mergeRange[T RangeEntryType, R Range[T]](new R, ranges []R) []R {
	out := make([]R, len(ranges), len(ranges)+1)
	copy(out, ranges)
	ranges = out
	for {
		mergerHappened := false
		for idx := range len(ranges) {
			intersection, hasIntersected := intersectRanges(new, ranges[idx])
			if hasIntersected {
				mergerHappened = true
				ranges = append(ranges[:idx], ranges[idx+1:]...)
				new = intersection
				break
			}
		}
		if !mergerHappened {
			ranges = append(ranges, new)
			break
		}
	}
	return ranges
}
