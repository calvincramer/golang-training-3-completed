package main

import (
	"io"
	"net/http"
	"time"
)

func IsPrime(n int64) bool {
	// Credit: http://stackoverflow.com/questions/1801391/what-is-the-best-algorithm-for-checking-if-a-number-is-prime
	switch {
	case n < 2:
		return false
	case n == 2:
		return true
	case n == 3:
		return true
	case n%2 == 0:
		return false
	case n%3 == 0:
		return false
	}

	var i int64 = 5
	var w int64 = 2

	for i*i <= n {
		if n%i == 0 {
			return false
		}
		i += w
		w = 6 - w
	}
	return true
}

// Do not modify me!
func GetGoogleWebpage() {
	// fmt.Println("Getting google.com...")
	time.Sleep(time.Millisecond * 750)
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Response has ", len(body), " bytes")
}

// Is prime result is given on the `res` channel
func IsPrimeGoroutine(n int64, res chan<- bool) {
	res <- IsPrime(n)
}

type SentinelT struct{}

var Sentinel SentinelT = SentinelT{}

type Command struct {
	Op   Operation
	Args []int
}

type Operation int

const (
	OP_SQUARE Operation = iota // Square the current value. No arguments accepted.
	OP_SET                     // Set the current value to Args[0]. One argument supplied.
	OP_ADD                     // Add each member of Args to the current value.
	OP_MULT                    // Multiple the current value by each member of Args
)

type Transaction struct {
	Category string
	Amount   int
}
