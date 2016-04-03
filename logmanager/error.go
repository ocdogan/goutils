package logmanager

import (
    "fmt"
    "runtime"
)

// Error struct for wrapping error with stack trace
type Error struct {
    err error
    stack string
}

// WrapErrorWithStack wraps the error with stack trace
func WrapErrorWithStack(e error) *Error {
    if e != nil {
        switch e.(type) {
            case *Error:
                return e.(*Error)
            default:
                stack := make([]byte, 1<<20)
                len := runtime.Stack(stack, true)
                
                return &Error {
                    err: e,
                    stack: string(stack[:len]),
                }
        }
    }
    return nil
}

// Error function used to implement the error interface
func (e *Error) Error() string {
    if e != nil && e.err != nil {
        return e.err.Error()
    }
    return ""
}

// String function used to implement the string interface
func (e *Error) String() string {
    if e != nil && e.err != nil {
        return fmt.Sprintf("Error: %s\nStack: %s\n", e.err.Error(), e.stack)
    }
    return ""
}