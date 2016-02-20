package stopwatch

import (
	"fmt"
	"time"
	"testing"
)

func debugLog(debug bool, str string, a ...interface{}) {
    if debug {
        fmt.Printf(str, a...)
    }
}

func mask1(s string, repeatcnt int, debug bool) string {
    debugLog(debug, "mask1, repeatcnt: %d\n", repeatcnt)
    
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
            debugLog(debug, "mask1, %d: %s\n", i, string(r))
            
            if i > 0 {
                if r == prev {
                    charcnt++
                    continue
                } 
                
                if charcnt >= repeatcnt {
                    stoppos = startpos+charcnt-1
                    
                    debugLog(debug, "mask1, startpos: %d, stoppos: %d\n", startpos, stoppos)
                    
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
                    
            debugLog(debug, "mask1, startpos: %d, stoppos: %d\n", startpos, stoppos)

            for j := stoppos; j >= startpos; j-- {
                charr[j] = '*'
            }
        }
        return string(charr)
    }
    return s
}

func mask2(s string, repeatcnt int, debug bool) string {
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
            debugLog(debug, "mask2, %d: %s\n", i, string(r))
                   
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
                    debugLog(debug, "mask2, startpos: %d\n", startpos)
                    
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

type masker func(s string, repeatcnt int, debug bool) string

func test(str string, debug bool, printResults bool, f masker) {
    s1 := f(str, 1, debug)
    s2 := f(str, 2, debug)
    s3 := f(str, 3, debug)
    s4 := f(str, 4, debug)
    s5 := f(str, 5, debug)
    s6 := f(str, 6, debug)
    
    if printResults {
        fmt.Println(0, str)
        fmt.Println(1, s1)
        fmt.Println(2, s2)
        fmt.Println(3, s3)
        fmt.Println(4, s4)
        fmt.Println(5, s5)
        fmt.Println(6, s6)
    }
}

func TestNano(t *testing.T) {
    str := "abbcccaaeeeeb bfffffca ccabbbb"
    
    w := New()
    for i := 0; i < 1000; i++ {
        test(str, false, false, mask2)
    }
    w.Stop()
    t2 := w.Milliseconds() * int64(time.Millisecond)

    if t2 <= 0 {
		t.Fatalf("Unexpected Value %d\n", t2)
	}

    w.Restart()
    for i := 0; i < 1000; i++ {
        test(str, false, false, mask1)
    }
    w.Stop()    
    t1 := w.Nanoseconds()

    if t1 <= 0 {
		t.Fatalf("Unexpected Value %d\n", t1)
	}

    fmt.Printf("* mask1, spent time: %d\n", t1)
    fmt.Printf("* mask2, spent time: %d\n", t2)    
}