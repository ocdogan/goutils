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
    "runtime"
    "time"
    "github.com/ocdogan/goutils/uuid"
)

type LogEntry struct {
    id string
    time time.Time
    duration time.Duration
    message string
    stack string
    logType LogType
    args map[string]interface{}
}

func NewLogEntry(message string, args map[string]interface{}) *LogEntry {
    uuid, _ := uuid.NewUUID()
    result := &LogEntry{
        id: uuid.String(),
        time: time.Now(),
        message: message,
        args: args,
    }
    
    if Enabled() && StacktraceEnabled() {
        stack := make([]byte, 1<<20)
        len := runtime.Stack(stack, true)
                    
        result.stack = string(stack[:len])
    }
    return result
}

func (entry *LogEntry) StartWatch() {
    entry.time = time.Now()
}

func (entry *LogEntry) StopWatch() {
    entry.duration = time.Now().Sub(entry.time)
}
