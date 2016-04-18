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
