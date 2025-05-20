package main

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
