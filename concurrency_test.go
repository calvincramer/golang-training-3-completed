package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func approxEq[T float32 | float64](n1 T, n2 T) bool {
	diff := n2 - n1
	if diff < 0 {
		diff = -diff
	}
	return diff < 1e-6
}

func TestDangerousSqrt1(t *testing.T) {
	if approxEq(DangerousSqrt(9.0), 3.0) != true {
		t.Fatalf("sqrt 9 should be 3")
	}
}

func TestDangerousSqrt2(t *testing.T) {
	defer func() {
		panicVal := recover()
		if panicVal == nil {
			t.Fatalf("Expected a panic, but none happened")
		} else {
			// Good, got panic.
		}
	}()
	_ = DangerousSqrt(-5)
}

func TestBubbleWrapSqrt(t *testing.T) {
	if approxEq(BubbleWrapSqrt(9.0), 3.0) != true {
		t.Fatalf("safe sqrt 9 should be 3")
	}
	if approxEq(BubbleWrapSqrt(-16.0), 4.0) != true {
		t.Fatalf("safe sqrt -16 should be 4")
	}
}

func TestIsPrimeMultiple(t *testing.T) {
	const REPEAT_TIMES = 100_000
	var parallelTime time.Duration
	var mainThreadTime time.Duration

	makePrimeNumbersSlice := func(repeat int) []int64 {
		const p1 int64 = 2147483423 // prime
		const p2 int64 = 2147482877 // prime
		const c1 int64 = 2147482879 // composite -> 227 x 1201 x 7877
		numbers := []int64{}
		for i := 1; i <= repeat; i++ {
			numbers = append(numbers, p1, p2, c1)
		}
		return numbers
	}

	checkResults := func(numbers []int64, res []bool) {
		for idx, numIsPrime := range res {
			switch idx % 3 {
			case 0, 1:
				if numIsPrime != true {
					t.Fatalf("Expected %d to be prime", numbers[idx])
				}
			case 2:
				if numIsPrime != false {
					t.Fatalf("Expected %d to be composite", numbers[idx])
				}
			}
		}
	}

	// Run multithreaded, check answers
	{
		numbers := makePrimeNumbersSlice(REPEAT_TIMES)
		start := time.Now()
		res := IsPrimeMultiple(numbers)
		parallelTime = time.Since(start)
		checkResults(numbers, res)
	}

	if runtime.NumCPU() <= 1 {
		t.Skipf("Cannot determine if IsPrimeMultiple spawns goroutines to accomplish work since we compare elapsed time using single thread, and the system has one core.")
	}

	// Run smaller part on main thread to ensure IsPrimeMultiple is using goroutines
	{
		divFactor := 20
		numbers := makePrimeNumbersSlice(REPEAT_TIMES / divFactor)
		start := time.Now()
		var result []bool = make([]bool, len(numbers))
		for idx, num := range numbers {
			result[idx] = IsPrime(num)
		}
		mainThreadTime = time.Since(start)
		checkResults(numbers, result)

		// Compare times. Parallel should be at least 50% faster than running on main thread.
		// Even a CPU with two cores should run close to 2x faster.
		mainThreadWholeSecs := mainThreadTime.Seconds() * float64(divFactor)
		parallelSecs := parallelTime.Seconds()
		speedup := mainThreadWholeSecs / parallelSecs

		fmt.Printf("%f - %f (%f) - %f\n", parallelSecs, mainThreadWholeSecs, mainThreadTime.Seconds(), speedup)

		if speedup < 1.5 {
			t.Fatalf("It looks like IsPrimeMultiple is not spawning goroutines based on execution time")
		}
	}
}
