package main

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

type V3 [3]int
type Circuit uint
type BoxPair [2]uint
type BoxPairSet map[BoxPair]struct{}

var V3Zero = V3{0, 0, 0}

func parse3DCoordinatesFromRow(line string) (V3, error) {
	parts := strings.Split(line, ",")
	values := [3]int{0, 0, 0}
	if len(parts) != 3 {
		return V3Zero, fmt.Errorf("invalid line: `%s`", line)
	}
	for idx := range 3 {
		val, err := strconv.Atoi(parts[idx])
		if err != nil {
			return V3Zero, fmt.Errorf("invalid coordinates in line `%s`", line)
		}
		values[idx] = val
	}
	return V3{values[0], values[1], values[2]}, nil
}

func v3Dist(fst, scn V3) uint {
	out := uint(0)
	for idx := range 3 {
		delta := fst[idx] - scn[idx]
		out += uint(delta * delta)
	}
	return out
}

func parse3DCoordinatesFromStdin(yield func(V3) bool) {
	for line := range parseLinesFromStdin {
		if line == "" {
			break
		}
		vec3d, err := parse3DCoordinatesFromRow(line)
		if err != nil {
			panic(fmt.Sprintf("invalid row `%s`", line))
		}
		if !yield(vec3d) {
			break
		}
	}
}

func closestV3(lst []V3, idx int) uint {
	fst := lst[idx]
	min_dist := -1
	min_idx := -1
	for scn_idx := range len(lst) {
		if idx == scn_idx {
			continue
		}
		cur_dist := int(v3Dist(fst, lst[scn_idx]))
		if min_dist >= 0 && cur_dist < min_dist {
			min_dist = int(cur_dist)
			min_idx = scn_idx
		}
	}
	return uint(min_idx)
}

func findWorst(values []uint) int {
	worst := uint(0)
	worst_idx := -1
	for idx, val := range values {
		if val > worst {
			worst = val
			worst_idx = idx
		}
	}
	return worst_idx
}

func makeClosestPairs(points []V3, n uint) []BoxPair {
	out := make([]BoxPair, n)
	distances := make([]uint, n)
	worst_distance := -1
	worst_index := 0
	added := uint(0)
	for fst := range len(points) {
		first_point := points[fst]
		for scn := range len(points) - fst - 1 {
			scn = fst + scn + 1
			second_point := points[scn]
			dist := v3Dist(first_point, second_point)
			if added < n {
				out[added] = BoxPair{uint(fst), uint(scn)}
				distances[added] = dist
				if worst_distance < 0 || dist > uint(worst_distance) {
					worst_distance = int(dist)
					worst_index = scn
				}
			} else {
				if dist > uint(worst_distance) {
					continue
				}
				distances[worst_index] = dist
				out[worst_index] = BoxPair{uint(fst), uint(scn)}
				worst_index = findWorst(distances)
				worst_distance = int(distances[worst_index])
			}
			added++
		}
	}
	return out
}

type BoxPairDistance struct {
	pair     BoxPair
	distance uint
}
type BoxPairHeap []BoxPairDistance

func (hp BoxPairHeap) Len() int           { return len(hp) }
func (hp BoxPairHeap) Less(i, j int) bool { return hp[i].distance > hp[j].distance }
func (hp BoxPairHeap) Swap(i, j int) {
	hp[i], hp[j] = hp[j], hp[i]
}
func (hp *BoxPairHeap) Push(x any) {
	item := x.(BoxPairDistance)
	*hp = append(*hp, item)
}
func (hp *BoxPairHeap) Pop() any {
	old := *hp
	n := len(old)
	out := old[n-1]
	(*hp) = old[:n-1]
	return out
}

func makeClosestPairsHeap(points []V3, n uint) []BoxPair {
	out := make([]BoxPair, n)
	hp := make(BoxPairHeap, 0)
	added := uint(0)
	for fst := range len(points) {
		first_point := points[fst]
		for scn := range len(points) - fst - 1 {
			scn = fst + scn + 1
			second_point := points[scn]
			dist := v3Dist(first_point, second_point)
			if added < n {
				heap.Push(&hp, BoxPairDistance{BoxPair{uint(fst), uint(scn)}, dist})
			} else {
				if dist >= hp[0].distance {
					continue
				}

				hp[0].distance = dist
				hp[0].pair = BoxPair{uint(fst), uint(scn)}
				heap.Fix(&hp, 0)
			}
			added++
		}
	}
	for idx := range n {
		out[n-idx-1] = heap.Pop(&hp).(BoxPairDistance).pair
	}
	return out
}

func mergeBoxes(i, j int, circuits []Circuit) {
	merge_to := min(circuits[i], circuits[j])
	merge_from := max(circuits[i], circuits[j])
	for pos, val := range circuits {
		if val == merge_from {
			circuits[pos] = merge_to
		}
	}
}

type CircuitCount struct {
	circuit Circuit
	count   uint
}
type CircuitHeap []CircuitCount

func (hp CircuitHeap) Len() int           { return len(hp) }
func (hp CircuitHeap) Less(i, j int) bool { return hp[i].count > hp[j].count }
func (hp CircuitHeap) Swap(i, j int) {
	hp[i], hp[j] = hp[j], hp[i]
}
func (hp *CircuitHeap) Push(x any) {
	item := x.(CircuitCount)
	*hp = append(*hp, item)
}
func (hp *CircuitHeap) Pop() any {
	old := *hp
	n := len(old)
	out := old[n-1]
	(*hp) = old[:n-1]
	return out
}

func takeTopCircuits(circuits []Circuit, n uint) []CircuitCount {
	counts := make([]uint, len(circuits))
	out := make([]CircuitCount, n)
	for idx := range len(counts) {
		counts[idx] = 0
	}
	for idx := range len(circuits) {
		counts[circuits[idx]]++
	}
	hp := make(CircuitHeap, 0)
	for idx, val := range counts {
		heap.Push(&hp, CircuitCount{Circuit(idx), val})
	}
	for idx := range n {
		out[idx] = heap.Pop(&hp).(CircuitCount)
	}
	return out
}

func aoc8a(n_distances uint, n_circuits uint) uint {
	password := uint(1)
	circuits := make([]Circuit, 0)
	points := make([]V3, 0)
	idx := 0

	for coord := range parse3DCoordinatesFromStdin {
		points = append(points, coord)
		circuits = append(circuits, Circuit(idx))
		idx++
	}
	closest_pairs := makeClosestPairsHeap(points, n_distances)
	for _, pair := range closest_pairs {
		mergeBoxes(int(pair[0]), int(pair[1]), circuits)
	}
	top_circuits := takeTopCircuits(circuits, n_circuits)
	for _, circuit := range top_circuits {
		fmt.Printf("circuit %d: %d\n", circuit.circuit, circuit.count)
		password *= circuit.count
	}
	return password
}
