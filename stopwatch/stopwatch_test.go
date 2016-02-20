//	The MIT License (MIT)
//
//	Copyright (c) 2015, Cagatay Dogan
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

package stopwatch

import (
	"fmt"
	"time"
	"testing"
)

var debug = false

func setDebug(state bool) {
    debug = state
}

func debugLog(str string, a ...interface{}) {
    if debug {
        fmt.Printf(str, a...)
    }
}

func mask1(s string, repeatcnt int) string {
    debugLog("mask1, repeatcnt: %d\n", repeatcnt)
    
    if repeatcnt < 1 {
        panic("Invalid value for param: repeatcnt")
    }
    
    if s != "" {
        charr := make([]rune, len(s))
        if repeatcnt == 1 {
            for i := 0; i < len(charr); i++ {
                charr[i] = '*'
            }
            return string(charr)
        }
        
        charcnt := 1
        startpos := 0
        var prev rune
        var stoppos int 
        
        for i, r := range s {
            charr[i] = r
            debugLog("mask1, %d: %s\n", i, string(r))
            
            if i > 0 {
                if r == prev {
                    charcnt++
                    continue
                } 
                
                if charcnt >= repeatcnt {
                    stoppos = startpos+charcnt-1
                    
                    debugLog("mask1, startpos: %d, stoppos: %d\n", startpos, stoppos)
                    
                    for j := stoppos; j >= startpos; j-- {
                        charr[j] = '*'
                    }
                }
                    
                charcnt = 1
                startpos = i
            }
            prev = r
        }
        
        if charcnt >= repeatcnt && startpos <= len(charr)-repeatcnt {
            stoppos = len(charr)-1
                    
            debugLog("mask1, startpos: %d, stoppos: %d\n", startpos, stoppos)

            for j := stoppos; j >= startpos; j-- {
                charr[j] = '*'
            }
        }
        return string(charr)
    }
    return s
}

func mask2(s string, repeatcnt int) string {
    if repeatcnt < 1 {
        panic("Invalid value for param: repeatcnt")
    }
    
    if s != "" {
        charr := make([]rune, len(s))
        if repeatcnt == 1 {
            for i := 0; i < len(charr); i++ {
                charr[i] = '*'
            }
            return string(charr)
        }
        
        var charcnt int
        var startpos int
        var prev rune
        var mask bool
        
        for i, r := range s {
            debugLog("mask2, %d: %s\n", i, string(r))
                   
            if i == 0 || r != prev {
                prev = r
                charcnt = 1
                startpos = i
                mask = false
                charr[i] = r
                
                continue
            }
            
            charcnt++
            if mask {
                charr[i] = '*'
                startpos = i + 1
            } else {
                mask = charcnt >= repeatcnt
                if !mask {
                    charr[i] = r
                } else {
                    debugLog("mask2, startpos: %d\n", startpos)
                    
                    for j := i; j >= startpos; j-- {
                        charr[j] = '*'
                    }
                    startpos = i + 1
                }
            }
        }
        return string(charr)
    }
    return s
}

type masker func(s string, repeatcnt int) string

func runTest(f masker, str string, printResults bool) {
    results := make([]string, 6)
    for i := range results {
        results[i] = f(str, i + 1)
    }
    
    if printResults {
        fmt.Println(0, str)
        for i, s := range results {
            fmt.Println(i + 1, s)
        }
    }
}

func TestNano(t *testing.T) {
    str := "abbcccaaeeeeb bfffffca ccabbbb"
    
    w := New()
    for i := 0; i < 1000; i++ {
        runTest(mask2, str, false)
    }
    w.Stop()
    t2 := w.Milliseconds() * int64(time.Millisecond)

    if t2 <= 0 {
		t.Fatalf("Unexpected Value %d\n", t2)
	}

    w.Restart()
    for i := 0; i < 1000; i++ {
        runTest(mask1, str, false)
    }
    w.Stop()    
    t1 := w.Nanoseconds()

    if t1 <= 0 {
		t.Fatalf("Unexpected Value %d\n", t1)
	}

    fmt.Printf("* mask1, spent time: %d\n", t1)
    fmt.Printf("* mask2, spent time: %d\n", t2)    
}