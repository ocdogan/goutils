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
