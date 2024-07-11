package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Make some random keys
func MakeKeys(num, length int) [][]int {
	ret := make([][]int, num)
	for i, _ := range ret {
		key := make([]int, length)
		for j, _ := range key {
			key[j] = rand.Intn(4)
		}
		ret[i] = key
	}
	return ret
}

// Return the number of matches between a and b
func compare(a, b []int) int {
	var same int
	for i, x := range a {
		if x == b[i] {
			same++
		}
	}
	return same
}

/*
Make two random streams (the 'A' stream and the 'B' stream) with a given
similarity and then look at what the similarity is in the regions where one or
both of them matches one of some preset random keys. If 'either' is true, the
regions are defined by either the A stream or the B stream matching one of the
keys. Otherwise they're defined just by the A stream matching. Return the
average similarity inside the regions for all matches.
*/
func CompareInRegions(streamLen int, similarity float64, either bool) float64 {
	var total, totalSame int

	// 8 keys of length 6 works well. If you make individual keys too long you
	// need a lot of iterations before much happens.
	keyLen := 6
	numKeys := 4
	keys := MakeKeys(numKeys, keyLen)

	// Keep the last keyLen outputs from each stream around so that we can
	// check if they match the keys.
	aBuf := make([]int, keyLen)
	bBuf := make([]int, keyLen)

	for i := 0; i < streamLen; i++ {
		// Generate the next character of the two streams.
		var a, b int
		a = rand.Intn(4)
		if rand.Float64() < similarity {
			b = a
		} else {
			for {
				b = rand.Intn(4)
				if b != a {
					break
				}
			}
		}

		aBuf[i%keyLen] = a
		bBuf[i%keyLen] = b

		// Is this part of the stream inside one of the regions defined by the
		// keys?
		insideRegions := false
		for _, key := range keys {
			insideRegions = reflect.DeepEqual(aBuf, key)
			if either {
				insideRegions = insideRegions || reflect.DeepEqual(bBuf, key)
			}
			if insideRegions {
				break
			}
		}

		if insideRegions {
			same := compare(aBuf, bBuf)
			totalSame += same
			total += len(aBuf)
		}
	}

	ret := float64(totalSame) / float64(total)
	return ret
}

func main() {
	var lower int
	trials := 100

	for i := 0; i < trials; i++ {
		t := CompareInRegions(1000000, 0.65, true)
		f := CompareInRegions(1000000, 0.65, false)
		if t < f {
			lower++
		}
		fmt.Printf("%d/%d completed\n", i+1, trials)
	}

	fmt.Printf("Similarity in the regions was lower with "+
		"either=true %d/%d times\n", lower, trials)
}
