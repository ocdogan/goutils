package logmanager

import (
    "sync"
    "sync/atomic"
)

var (
    isEnabled = enabled
    isStacktraceEnabled = disabled
    bucketMtx = &sync.Mutex{}
    buckets = make(map[string]*logBucket)
)

// Enable enables the logging manager globally
func Enable() {
    atomic.StoreUint32(&isEnabled, enabled)
}

// Disable disables the logging manager globally
func Disable() {
    atomic.StoreUint32(&isEnabled, disabled)
}

// Enabled function is used to get if the global logging manager is enabled
func Enabled() bool {
    return atomic.LoadUint32(&isEnabled) == enabled
}

// EnableStacktrace enables using stacktrace in logging
func EnableStacktrace() {
    atomic.StoreUint32(&isStacktraceEnabled, enabled)
}

// DisableStacktrace disables using stacktrace in logging
func DisableStacktrace() {
    atomic.StoreUint32(&isStacktraceEnabled, disabled)
}

// StacktraceEnabled function is used to get if the global stacktrace usage is enabled
func StacktraceEnabled() bool {
    return atomic.LoadUint32(&isStacktraceEnabled) == enabled
}

// RegisterHandler adds the handler into the logging chain
func RegisterHandler(handler LogHandler) {
    if handler != nil {
        name := handler.Name()
        
        UnregisterHandler(name)        
        bucketMtx.Lock()
        defer bucketMtx.Unlock()
        
        bucket := newBucket(handler)
        buckets[handler.Name()] = bucket
        go bucket.process()
    }
}

// UnregisterHandler removes the handler from the logging chain
func UnregisterHandler(name string) {
    bucketMtx.Lock()
    defer bucketMtx.Unlock()
    
    bucket, ok := buckets[name]
    if ok {
        bucket.close()
    }
    delete(buckets, name)
}

// LogError is used to log the given error by log manager
func LogError(e error, args map[string]interface{}) {
    if e != nil && Enabled() {
        entry := NewLogEntry(e.Error(), args)
        entry.logType = ErrorLog
        Log(entry)
    }
}

// LogMessage is used to log the given message by log manager
func LogMessage(message string, args map[string]interface{}) {
    if Enabled() {
        Log(NewLogEntry(message, args))
    }    
}

// Log lets the given entry to be processes by the handler chain
func Log(entry *LogEntry) {
    if entry != nil && Enabled() {
        bucketMtx.Lock()
        defer bucketMtx.Unlock()
        
        for _, bucket := range buckets {
            if bucket.enabled() {
                bucket.entryChan <- entry
            }
        }
    }    
}
