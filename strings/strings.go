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

package strings

import (
	"bytes"
	"io"
	str "strings"
	"unicode"
)

// EndsWith returns if s ends with end.
func EndsWith(s, end string) bool {
	lenS := len(s)
	if lenS == 0 {
		return len(end) == 0
	}

	lenEnd := len(end)
	if lenEnd == 0 || lenEnd > lenS {
		return false
	}

	return s[lenS-lenEnd:lenS] == end
}

// StartsWith returns if s starts with start.
func StartsWith(s, start string) bool {
	lenS := len(s)
	if lenS == 0 {
		return len(start) == 0
	}

	lenStart := len(start)
	if lenStart == 0 || lenStart > lenS {
		return false
	}

	return s[:lenStart] == start
}

// Capitalize returns the capitalized form of the given string.
func Capitalize(s string) string {
	if len(s) > 0 {
		inBuf := bytes.NewBufferString(s)
		outBuf := bytes.NewBufferString("")

		space := false
		punct := false

		for i := 0; ; i++ {
			rn, size, e := inBuf.ReadRune()
			if size == 0 || e == io.EOF {
				break
			}
			if e != nil {
				return s
			}

			if unicode.IsSpace(rn) {
				space = true
				punct = false
			} else if unicode.IsPunct(rn) {
				punct = true
				space = false
			} else if space || punct || i == 0 {
				space = false
				punct = false
				if unicode.IsLetter(rn) {
					rn = unicode.ToUpper(rn)
				}
			}

			outBuf.WriteRune(rn)
		}
		return outBuf.String()
	}
	return s
}

// TrimLeft returns a slice of the string s, with all leading
// and trailing white space removed from left, as defined by Unicode.
func TrimLeft(s string) string {
	if len(s) > 0 {
		return str.TrimLeftFunc(s, unicode.IsSpace)
	}
	return s
}

// TrimRight returns a slice of the string s, with all leading
// and trailing white space removed from right, as defined by Unicode.
func TrimRight(s string) string {
	if len(s) > 0 {
		return str.TrimRightFunc(s, unicode.IsSpace)
	}
	return s
}
