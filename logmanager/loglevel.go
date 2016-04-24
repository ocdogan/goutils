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
    "bytes"
)

// LogLevel is used to inform the system about the given message type
type LogLevel byte

const (
    // LevelInfo is used if the message is an information
    LevelInfo LogLevel = 1
    // LevelWarning is used if the message is a warning
    LevelWarning LogLevel = 2
    // LevelError is used if the message is an error
    LevelError LogLevel = 4
    // LevelFatal is used if the message is a fatal error
    LevelFatal LogLevel = 8
    // AllLogLevels is used to push a logentry to a handler with overriding its LoggingType
    AllLogLevels = LevelInfo | LevelWarning | LevelError | LevelFatal
)

var (
    allLogLevels = []LogLevel{ LevelInfo, LevelWarning, LevelError, LevelFatal }
)

// Has checks if lt includes lt2 
func (l LogLevel) Has(l2 LogLevel) bool {
    return l & l2 == l2
}

func (l LogLevel) toString() string {
    switch l {
    case LevelInfo:
        return "info"
    case LevelWarning:
        return "warning"
    case LevelError:
        return "error"
    case LevelFatal:
        return "fatal"
    }
    return ""
}

func (l LogLevel) String() string {
    if l == LogLevel(0) {
        l = AllLogLevels
    }
    
    buf := &bytes.Buffer{}
    for _, l2 := range allLogLevels {
        if l.Has(l2) {
            if buf.Len() > 0 {
                buf.WriteString("|")
            }
            buf.WriteString(l2.toString())
        }
    }
    return buf.String()
}