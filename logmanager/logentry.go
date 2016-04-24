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
    "encoding/json"
    "fmt"
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

// ID returns the id of the entry
func (entry *LogEntry) ID() string {
    return entry.id
}

// Time returns the creation time of the entry
func (entry *LogEntry) Time() time.Time {
    return entry.time
}

// Duration returns the measured duration (time passed between StartWatch and StopWatch) of the entry
func (entry *LogEntry) Duration() time.Duration {
    return entry.duration
}

// Message returns the message of the entry
func (entry *LogEntry) Message() string {
    return entry.message
}

// Stack returns the stack trace of the entry
func (entry *LogEntry) Stack() string {
    return entry.stack
}

// LogType returns the log type of the entry
func (entry *LogEntry) LogType() LogType {
    return entry.logType
}

// Args returns the arguments of the entry
func (entry *LogEntry) Args() map[string]interface{} {
    return entry.args
}

// StartWatch starts timer to measure the time passed
func (entry *LogEntry) StartWatch() {
    entry.time = time.Now()
}

// StopWatch stops the timer to measure the time passed
func (entry *LogEntry) StopWatch() {
    entry.duration = time.Now().Sub(entry.time)
}

// ToJSON returns the JSON formatted entry as byte array
func (entry *LogEntry) ToJSON() []byte {
    b, e := json.Marshal(struct{
        ID string `json:"id"`
        Time time.Time `json:"time"`
        Duration time.Duration `json:"duration"`
        LogType string `json:"logType"`
        Message string `json:"message"`
        Stack string `json:"stack"`
        Args map[string]interface{} `json:"args,omitempty"`
    }{
        ID: entry.id,
        Time: entry.time,
        Duration: entry.duration,
        LogType: entry.logType.String(),
        Message: entry.message,
        Stack: entry.stack,
        Args: entry.args,
    })
    if e != nil {
        return nil
    }
    return b
}

// ToJSON returns the text line formatted entry as byte array
func (entry *LogEntry) ToText() []byte {
    if entry == nil {
        return nil
    }
    
    buffer := &bytes.Buffer{}

    writeKvToBuffer(buffer, "id", entry.id)
    writeKvToBuffer(buffer, "time", entry.time)
    writeKvToBuffer(buffer, "duration", entry.duration)
    writeKvToBuffer(buffer, "logType", entry.logType.String())
    writeKvToBuffer(buffer, "message", entry.message)
    writeKvToBuffer(buffer, "stack", entry.stack)
    
    if entry.args != nil {
        for k, v := range entry.args {
            writeKvToBuffer(buffer, k, v)
        }
    }
    
    return buffer.Bytes()
}

func textOrNumber(value string) bool {
	for _, c := range value {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
            c == '-' || c == '+' || 
            c == '.' || c == ',') {
			return false
		}
	}
	return true
}

func writeKvToBuffer(buffer *bytes.Buffer, key string, value interface{}) {
    buffer.WriteString(key)
    
    switch value.(type) {
    case nil:
        break
    case string:
        if !textOrNumber(value.(string)) {
            fmt.Fprintf(buffer, "=%q ", value)
        } else {
            buffer.WriteString("=\"")
            buffer.WriteString(value.(string))
            buffer.WriteString("\" ")
        }
    default:
        buffer.WriteString("=\"")
        fmt.Fprint(buffer, value)         
        buffer.WriteString("\" ")
    }
}