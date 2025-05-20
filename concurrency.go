package main

import (
	"math"
	"sync"
)

// TODO: write a function that takes the square root of a number.
// If the number is negative, then panic (we don't know what imaginary numbers are).
func DangerousSqrt(n float64) float64 {
	if n < 0 {
		panic("Square root negative number is impossible!")
	}
	return math.Sqrt(n)
}

// TODO: write a function that uses DangerousSqrt() to accomplish a safe square root. In the case
// that DangerousSqrt() panics, recover and try DangerousSqrt() again with the absolute value of n.
func BubbleWrapSqrt(n float64) (ans float64) {
	defer func() {
		panicVal := recover()
		if panicVal != nil {
			ans = DangerousSqrt(-n)
		}
	}()
	ans = DangerousSqrt(n)
	return
}

// TODO: simple goroutine do something in background, wait for it to finish.
func DoBackground() {
	// TODO CALVIN!!!
}

// TODO: spawn a goroutine for each number to check if it is prime or not. This method will be
// called with a large amount of very large numbers. The starting implementation calls IsPrime()
// blocking on the current thread.
//
// Hint: use IsPrime() which is defined in util.go
// Hint: use sync.WaitGroup to wait for multiple goroutines
// Hint: you should see all cores of your CPU be saturated for a few seconds. Use 'top', 'htop', or
// your system monitor to see. You should **NOT** see a single core at 100%.
//
// Note: this function should take around 1 second to compute. Adjust `REPEAT_TIMES` in
// concurrency_test.go if your machine is slower or faster.
// Note: there is no synchronization needed to modify the result slice because no two goroutines
// will modify the same index.
func IsPrimeMultiple(numbers []int64) []bool {
	var waitGroup sync.WaitGroup
	var result []bool = make([]bool, len(numbers))
	for idx, num := range numbers {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			result[idx] = IsPrime(num)
		}()
	}
	waitGroup.Wait()
	return result
}

// TODO CALVIN: unbuffered channel example
// TODO CALVIN: buffered channel example
// TODO CALVIN: select on multiple channels
// TODO CALVIN: synchronization primitive example, multiple goroutines writing same value
