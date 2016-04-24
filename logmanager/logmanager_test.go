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

package logmanager

import (
    "fmt"
    "runtime"
    "testing"
    "time"
)

func TestLogManager(t *testing.T) {
    fmt.Println("\nTestLogManager\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

    mem1 := new(runtime.MemStats)
    runtime.ReadMemStats(mem1)
    
    RegisterHandler(&ConsoleLogHandler{})
    RegisterHandlerWithName("x", &ConsoleLogHandler{})

    args := map[string]interface{}{
        "a": "1",
        "b": 2,
        "c": 3.3,
        "d": map[string]interface{}{
            "e": "x",
            "f": "y",
        },
    }
    for i := 0; i < 20; i++ {
        LogMessage("xxx", args)
    }
    time.Sleep(10*time.Second)
    
    mem2 := new(runtime.MemStats)
    runtime.ReadMemStats(mem2)
    if mem2.Alloc <= mem1.Alloc {
        fmt.Printf("Mem allocated: 0 MB\n")
    } else {
        fmt.Printf("Mem allocated: %3.3f MB\n", float64(mem2.Alloc - mem1.Alloc)/(1024*1024))
    }
}
