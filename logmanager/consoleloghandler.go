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
    "sync/atomic"
)

// ConsoleLogHandler is used to push the log entry to console
type ConsoleLogHandler struct {
    disabled uint32
}

// Name returns the name of the handler used for registration
func (handler *ConsoleLogHandler) Name() string {
    return "console"
}

// Enabled returns if the handler is active
func (handler *ConsoleLogHandler) Enabled() bool {
    return atomic.LoadUint32(&handler.disabled) == falseUint32
}

// Enable activates the handler
func (handler *ConsoleLogHandler) Enable() {
    atomic.StoreUint32(&handler.disabled, falseUint32)
}

// Disable deactivates the handler
func (handler *ConsoleLogHandler) Disable() {
    atomic.StoreUint32(&handler.disabled, trueUint32)
}

// Level gives if the pushed entry should be logged by the handler
func (handler *ConsoleLogHandler) Level() LogLevel {
    return AllLogLevels
}

// Format gives the format that will be used by the handler
func (handler *ConsoleLogHandler) Format() LogFormat {
    return JSONFormat
}

// QueueLen gives the queue length that will be used when the entry is queued
func (handler *ConsoleLogHandler) QueueLen() int {
    return -1
}

// Process evaluates the given entry
func (handler *ConsoleLogHandler) Process(entry interface{}) {
    if data, ok := entry.([]byte); ok {
        fmt.Println(string(data))
    }
}

