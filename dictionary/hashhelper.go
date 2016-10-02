//	The MIT License (MIT)
//
//	Copyright (c) 2016, Cagatay Dogan
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//		The above copyright notice and this permission notice shall be included in
//		all copies or substantial portions of the Software.
//
//		THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//		IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//		FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//		AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//		LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//		OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//		THE SOFTWARE.

package dictionary

import (
	"errors"
	"math"
)

const hashPrime int = 101
const MaxPrimeArrayLength int = 0x7FEFFFFD

var primes = []int{
	3, 7, 11, 17, 23, 29, 37, 47, 59, 71, 89, 107, 131, 163, 197, 239, 293, 353, 431, 521, 631, 761, 919,
	1103, 1327, 1597, 1931, 2333, 2801, 3371, 4049, 4861, 5839, 7013, 8419, 10103, 12143, 14591,
	17519, 21023, 25229, 30293, 36353, 43627, 52361, 62851, 75431, 90523, 108631, 130363, 156437,
	187751, 225307, 270371, 324449, 389357, 467237, 560689, 672827, 807403, 968897, 1162687, 1395263,
	1674319, 2009191, 2411033, 2893249, 3471899, 4166287, 4999559, 5999471, 7199369}

func IsPrime(candidate int) bool {
	if (candidate & 1) != 0 {
		limit := candidate * candidate
		for divisor := 3; divisor <= limit; divisor += 2 {
			if (candidate % divisor) == 0 {
				return false
			}
		}
		return true
	}
	return (candidate == 2)
}

func GetPrime(min int) (int, error) {
	if min < 0 {
		return 0, errors.New("Capacity overflow")
	}

	primesLen := len(primes)
	for i := 0; i < primesLen; i++ {
		prime := primes[i]
		if prime >= min {
			return prime, nil
		}
	}

	for i := (min | 1); i < math.MaxInt32; i += 2 {
		if IsPrime(i) && ((i-1)%hashPrime != 0) {
			return i, nil
		}
	}
	return min, nil
}

func ExpandPrime(oldSize int) int {
	newSize := 2 * oldSize
	if uint(newSize) > uint(MaxPrimeArrayLength) && MaxPrimeArrayLength > oldSize {
		return MaxPrimeArrayLength
	}

	result, _ := GetPrime(newSize)
	return result
}
