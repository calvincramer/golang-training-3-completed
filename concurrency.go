package main

import (
	"math"
	"sync"
	"time"
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

// TODO: run the `GetGoogleWebpage()` function in the background using goroutines. Do NOT wait for
// the goroutine to finish. Do NOT modify the `GetGoogleWebpage()` function.
func DoBackground() {
	go GetGoogleWebpage() // Make me a background task!
}

// TODO: the follow method attempts to modify the `sum` global variable using multiple goroutines to
// exploit the systems multiple CPU cores. It incorrectly calculates the sum of number from 0 to one
// million. Fix it so it is correctly calculated.
// Hint: the problem is a race condition
// Hint: you can use a lock or channels
var sum uint64 = 0
var sumLock sync.Mutex

func CalcSumOneMillion() uint64 {
	worker := func(low uint64, high uint64) {
		var tempSum uint64 = 0
		for n := low; n <= high; n++ {
			tempSum += n
		}
		sumLock.Lock()
		defer sumLock.Unlock()
		sum += tempSum // update global
	}
	sum = 0 // reset global
	var chunk uint64 = 100
	const end uint64 = 1_000_000
	var n uint64 = 0
	for n <= end {
		high := min(n+chunk-1, end)
		go worker(n, high)
		n += chunk
	}
	time.Sleep(time.Millisecond * 200) // bad way to wait for goroutines to finish!!!
	return sum
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

// TODO: run the `IsPrimeGoroutine()` as a goroutine and collect and return the response from the
// result channel.
func IsPrimeBackground(n int64) bool {
	ch := make(chan bool)
	go IsPrimeGoroutine(n, ch)
	return <-ch
}

// TODO: use an unbuffered channel to "ping pong" a piece of data between two goroutines, and then
// return the response. Here's the order of operations: (Sally and Carl identify the two threads)
//
//  1. Sally: multiple number by 5
//  2. Carl: if the number is negative then take its absolute value and add 100
//  3. Sally: square the number
//  4. Carl: divide the number by it's least significant digit (skip if the digit is 0)
//  5. Sally: subtract 7
//  6. Carl OR Sally: return the number to the caller.
//
// Note: in between each operation the number is sent to the other thread.
// Hint: you can use the current goroutine as one of the goroutines. No need to create two
// additional routines.
func PingPongCalc(num int) int {
	carl := func(ch chan int) {
		// We are Carl
		num := <-ch
		// #2
		if num < 0 {
			num = (-num) + 100
		}
		ch <- num
		// #4
		num = <-ch
		digit := num % 10
		if digit != 0 {
			num /= digit
		}
		ch <- num
	}

	ch := make(chan int)
	go carl(ch)

	// We are Sally
	// #1
	num *= 5
	ch <- num
	// #3
	num = <-ch
	num *= num
	ch <- num
	// #5
	num = <-ch
	num -= 7
	return num
}

// TODO CALVIN: buffered channel example
// TODO CALVIN: select on multiple channels
// TODO CALVIN: synchronization primitive example, multiple goroutines writing same value
