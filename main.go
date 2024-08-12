package main

import (
	"fmt"
	"math/rand"
	"time"
)

func binomialDistribution(marbleCount int, bucketCount int) []int {

	g := marbleGenerator(marbleCount)
	r := nextRow([]chan bool{g})

	for i := 1; i < bucketCount-1; i++ {
		r = nextRow(r)
	}

	return buckets(r)

}

func marbleGenerator(count int) chan bool {
	out := make(chan bool)

	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			out <- true
		}
	}()
	return out
}

func nextRow(parents []chan bool) []chan bool {
	children := make([]chan bool, len(parents)+1)

	children[0] = make(chan bool)

	for i := 1; i <= len(parents); i++ {
		children[i] = make(chan bool)

		go func() {
			for range parents[i-1] {
				r := rand.Intn(2)
				if r == 0 {
					children[i-1] <- true
				} else {
					children[i] <- true
				}
			}
		}()

	}

	return children
}

func buckets(leafNodes []chan bool) []int {
	result := make([]int, len(leafNodes))

	for i := 0; i < len(leafNodes); i++ {
		go func() {
			for range leafNodes[i] {
				result[i]++
			}
		}()
	}
	return result
}

// Iterate through the final values, printing out a set of * characters for
// each value, proportional to the value. If the longest value is >
// 80, scale all values down to fit within 80 characters.
func printOutput(values []int) {

	max := 0
	scaled := 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	for _, v := range values {
		if max > 80 {
			scaled = int(float64(v) / float64(max) * 80)
		} else {
			scaled = v
		}

		// Print the unscaled value, right justified and using
		// a fixed width of the max value.
		fmt.Printf("%*d ", len(fmt.Sprint(max)), v)
		for i := 0; i < scaled; i++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

func main() {
	values := binomialDistribution(500, 5)
	time.Sleep(1 * time.Second)
	printOutput(values)
}
