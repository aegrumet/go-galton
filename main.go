package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Simulates a Galton board using marbleCount marbles and binCount bins. Returns
// a slice of integers and a wait group. The slice of integers contains the
// number of marbles in each bin when the simulation is complete.
func binomialDistribution(marbleCount int, binCount int) ([]int, *sync.WaitGroup) {

	var wg sync.WaitGroup

	g := marbleSource(marbleCount)
	r := nextRow([]chan bool{g})

	for i := 1; i < binCount-1; i++ {
		r = nextRow(r)
	}

	wg.Add(binCount)
	return bins(&wg, r), &wg

}

// Returns a bool channel that will emit count true values before closing.
func marbleSource(count int) chan bool {
	out := make(chan bool)

	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			out <- true
		}
	}()
	return out
}

// Given a parent slice of channels, returns the next row of channels. Each
// channel in the returned slice can receive a value from the left or right
// parent channel (except at the boundaries), and will emit a value to the
// next row.
func nextRow(parents []chan bool) []chan bool {
	children := make([]chan bool, len(parents)+1)
	waitGroups := make([]sync.WaitGroup, len(parents)+1)

	children[0] = make(chan bool)

	for i := 1; i <= len(parents); i++ {
		children[i] = make(chan bool)

		waitGroups[i-1].Add(1)
		waitGroups[i].Add(1)

		go func() {
			for range parents[i-1] {
				r := rand.Intn(2)
				if r == 0 {
					children[i-1] <- true
				} else {
					children[i] <- true
				}
			}
			waitGroups[i-1].Done()
			waitGroups[i].Done()
		}()

		go func() {
			waitGroups[i].Wait()
			close(children[i])
		}()
	}

	go func() {
		waitGroups[0].Wait()
		close(children[0])
	}()

	return children
}

// Count the total number of marbles in each bin. Returns a slice of integers.
func bins(wg *sync.WaitGroup, leafNodes []chan bool) []int {
	result := make([]int, len(leafNodes))

	for i := 0; i < len(leafNodes); i++ {
		go func() {
			for range leafNodes[i] {
				result[i]++
			}
			wg.Done()
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
	values, wg := binomialDistribution(1000, 5)
	wg.Wait()
	printOutput(values)
}
