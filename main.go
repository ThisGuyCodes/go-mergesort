package main

import (
	"flag"
	"log"
	"math/rand"
	sortpkg "sort"
)

var (
	size = flag.Int("size", (1024*1024/4)*256, "how big a slice to sort")
)

func main() {
	flag.Parse()

	unsorted := make([]int, *size)

	for i := range unsorted {
		unsorted[i] = rand.Int()
	}

	scratch := make([]int, len(unsorted))

	log.Println("starting")
	sorted := <-sort(unsorted, scratch)
	log.Println("done")
	log.Println(sorted[0], sorted[len(sorted)-1])
	log.Println(sortpkg.IntsAreSorted(sorted))

}

// sort sorts data[], using tmp[] as merge space
func sort(data, tmp []int) <-chan []int {
	if len(data) <= 2 {
		retChan := make(chan []int, 1)
		if len(data) == 2 {
			if data[0] > data[1] {
				data[0], data[1] = data[1], data[0]
			}
		}
		retChan <- data
		return retChan
	}

	var m0, m1, m2 int

	if len(data) >= 4 {
		m0 = len(data) / 4
		m1 = m0 * 2
		m2 = m0 * 3
	} else if len(data) == 3 {
		m0 = 1
		m1 = 2
		m2 = 3
	} else if len(data) == 2 {
		m0 = 1
		m1 = 2
		m2 = 2
	}

	left := sort(data[:m0], tmp[:m0])
	leftMiddle := sort(data[m0:m1], tmp[m0:m1])
	rightMiddle := sort(data[m1:m2], tmp[m1:m2])
	right := sort(data[m2:], tmp[m2:])

	leftChan := merge(left, leftMiddle, tmp[:m1])
	rightChan := merge(rightMiddle, right, tmp[m1:])

	return merge(leftChan, rightChan, data)
}

// merge merges left and right into out
func merge(leftChan, rightChan <-chan []int, out []int) <-chan []int {
	retChan := make(chan []int)

	go func() {
		left, right := <-leftChan, <-rightChan
		if len(right) > len(left) {
			left, right = right, left
		}

		k := 0
		for ; len(left) > 0 && len(right) > 0; k++ {
			if left[0] < right[0] {
				out[k] = left[0]
				left = left[1:]
			} else {
				out[k] = right[0]
				right = right[1:]
			}
		}

		if len(left) == 0 {
			copy(out[k:], right)
		} else {
			copy(out[k:], left)
		}

		retChan <- out
	}()
	return retChan
}
