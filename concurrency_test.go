package main

import (
	"math/rand/v2"
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

func TestDoBackground(t *testing.T) {
	start := time.Now()
	DoBackground()
	elapsed := time.Since(start)
	if elapsed >= time.Millisecond*500 {
		t.Fatalf("Task is not ran in background")
	}
}

func TestCalcSumOneMillion(t *testing.T) {
	const correctAns uint64 = (1_000_000 * (1_000_000 + 1)) / 2
	if ans := CalcSumOneMillion(); ans != correctAns {
		t.Fatalf("Incorrect answer: %d", ans)
	}
}

func TestIsPrimeMultiple(t *testing.T) {
	const REPEAT_TIMES = 10_000
	const divFactor = 20
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
		numbers := makePrimeNumbersSlice(REPEAT_TIMES / divFactor)
		start := time.Now()
		var result []bool = make([]bool, len(numbers))
		for idx, num := range numbers {
			result[idx] = IsPrime(num)
		}
		mainThreadTime = time.Since(start)
		checkResults(numbers, result)
	}

	// Compare times. Parallel should be at least 50% faster than running on main thread.
	// Even a CPU with two cores should run close to 2x faster.
	mainThreadWholeSecs := mainThreadTime.Seconds() * float64(divFactor)
	parallelSecs := parallelTime.Seconds()
	speedup := mainThreadWholeSecs / parallelSecs

	// fmt.Printf("parallel: %f\n", parallelSecs)
	// fmt.Printf("main estimated: %f (direct: %f)\n", mainThreadWholeSecs, mainThreadTime.Seconds())
	// fmt.Printf("speedup: %f\n", speedup)

	if speedup < 1.5 {
		t.Fatalf("It looks like IsPrimeMultiple is not spawning goroutines based on execution time")
	}
}

func TestIsPrimeBackground(t *testing.T) {
	// Not testing that goroutines are used
	for _, n := range []int64{2, 3, 5, 7, 11, 13, 17, 19, 97} {
		if IsPrimeBackground(n) != true {
			t.Fatalf("Expected %d to be prime", n)
		}
	}
	for _, n := range []int64{1, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20, 100} {
		if IsPrimeBackground(n) != false {
			t.Fatalf("Expected %d to be composite", n)
		}
	}
}

func TestPingPongCalc(t *testing.T) {
	if PingPongCalc(5) != 118 {
		t.Fatalf("Incorrect answer for 5")
	}
	if PingPongCalc(-8) != 19593 {
		t.Fatalf("Incorrect answer for -8")
	}
	if PingPongCalc(-11) != 4798 {
		t.Fatalf("Incorrect answer for -11")
	}
}

func TestMultiWriter(t *testing.T) {
	var ch chan any
	// Make sure calling user function does not block, immediately returns channel
	{
		start := time.Now()
		ch = MultiWriter()
		if elapsed := time.Since(start); elapsed > time.Millisecond*100 {
			t.Fatalf("Calling MultiWriter() took longer than 100 ms. MultiWriter() should return a channel immediately.")
		}
	}
	var numInts int = 0
	var numStrs int = 0
	var numSentinels int = 0
done:
	for {
		select {
		case <-time.After(time.Second):
			t.Fatalf("timeout after 1 second, not getting expected output from channel!")
		case val := <-ch:
			switch val.(type) {
			case int:
				numInts += 1
				if numInts > 500 {
					t.Fatalf("Got more than 500 integers on the channel")
				}
			case string:
				numStrs += 1
				if numStrs > 500 {
					t.Fatalf("Got more than 500 strings on the channel")
				}
			case SentinelT:
				numSentinels += 1
				if numSentinels > 2 {
					t.Fatalf("Got more than 2 sentinels on the channel")
				}
				if numSentinels == 1 {
					if numInts < 500 && numStrs < 500 {
						t.Fatalf("Got first sentinel but did not receive 500 integers or 500 strings yet")
					}
					if numInts == 500 && numStrs == 500 {
						break done
					}
				} else if numSentinels == 2 {
					if numInts != 500 && numStrs != 500 {
						t.Fatalf("Got second sentinel, but do not have 500 integers and 500 strings")
					}
					break done
				}
			default:
				t.Fatalf("Unexpected type on channel: %T", val)
			}
		}
	}
	if numInts != 500 || numStrs != 500 || (numSentinels != 1 && numSentinels != 2) {
		t.Fatalf("Did not receive 500 integers and 500 strings and 1 or 2 sentinels")
	}
}

func TestCommandConsumer(t *testing.T) {
	{
		cmdCh := make(chan Command)
		doneCh := make(chan bool)
		go func(cmdCh chan Command, doneCh chan bool) {
			cmdCh <- Command{Op: OP_ADD, Args: []int{5, 8}}
			doneCh <- true
		}(cmdCh, doneCh)
		ret := CommandConsumer(cmdCh, doneCh)
		if ret != 13 {
			t.Fatalf("Incorrect answer: %d, expected 13", ret)
		}
	}
	{
		cmdCh := make(chan Command)
		doneCh := make(chan bool)
		go func(cmdCh chan Command, doneCh chan bool) {
			cmdCh <- Command{Op: OP_ADD, Args: []int{5, 8, -20}}
			cmdCh <- Command{Op: OP_SET, Args: []int{4}}
			cmdCh <- Command{Op: OP_SQUARE, Args: []int{}}
			cmdCh <- Command{Op: OP_MULT, Args: []int{-1, 2}}
			doneCh <- true
		}(cmdCh, doneCh)
		ret := CommandConsumer(cmdCh, doneCh)
		if ret != -32 {
			t.Fatalf("Incorrect answer: %d, expected -32", ret)
		}
	}
}

func TestPerformAccounting(t *testing.T) {
	transactions := []Transaction{}
	for i := 1; i <= 100_000; i++ {
		transactions = append(transactions, Transaction{Category: "travel", Amount: -1})
		transactions = append(transactions, Transaction{Category: "travel", Amount: 2})
		transactions = append(transactions, Transaction{Category: "food", Amount: 15})
		transactions = append(transactions, Transaction{Category: "food", Amount: -5})
	}
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})
	PerformAccounting(transactions)
	if CategoryTotals["travel"] != 100_000 {
		t.Fatalf("Expected travel category to be sum to 500000, got %d", CategoryTotals["travel"])
	}
	if CategoryTotals["food"] != 1_000_000 {
		t.Fatalf("Expected travel category to be sum to 5000000, got %d", CategoryTotals["food"])
	}
}
