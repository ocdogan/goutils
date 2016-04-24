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

package murmur

import (
    "encoding/binary"
)

const (
    u0  uint32 = 0
    u1  uint32 = 1
    u2  uint32 = 2
    u3  uint32 = 3
    u5  uint32 = 5
    b13 byte   = 13
    b15 byte   = 15
    c1  uint32 = 0xcc9e2d51
    c2  uint32 = 0x1b873593
    k   uint32 = 0xe6546b64
    fx1 uint32 = 0x85ebca6b
    fx2 uint32 = 0xc2b2ae35
)

func MurmurHash3(data []byte, length uint32, seed uint32) uint32 {
    h1 := seed
    nblocks := length >> 2

    i := 0
    for j := nblocks; j > u0; j-- {
        k1l := binary.LittleEndian.Uint32(data[i:])

        k1l *= c1
        k1l = mmhRotateLeft(k1l, b15)
        k1l *= c2

        h1 ^= k1l
        h1 = mmhRotateLeft(h1, b13)
        h1 = h1 * u5 + k

        i += 4
    }

    k1 := u0
    nblocks <<= u2
    tailLength := length & u3

    if tailLength == u3 {
        k1 ^= uint32(data[2 + nblocks]) << 16
    }
    if tailLength >= u2 {
        k1 ^= uint32(data[1 + nblocks]) << 8
    }
    if tailLength >= u1 {
        k1 ^= uint32(data[nblocks])
        k1 *= c1 
        k1 = mmhRotateLeft(k1, b15) 
        k1 *= c2
        h1 ^= k1
    }

    h1 ^= length
    return mmhFMix(h1)
}

func mmhFMix(h uint32) uint32 {
    h ^= h >> 16
    h *= fx1
    h ^= h >> 13
    h *= fx2
    h ^= h >> 16

    return h
}

func mmhRotateLeft(x uint32, r byte) uint32 {
    return (x << r) | (x >> (32 - r))
}
