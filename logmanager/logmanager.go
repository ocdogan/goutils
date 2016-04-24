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
        buckets[name] = bucket
        go bucket.process()
    }
}

// RegisterHandlerWithName adds the handler into the logging chain with the given name
func RegisterHandlerWithName(name string, handler LogHandler) {
    if handler != nil {
        if name == "" {
            name = handler.Name()
        }
        
        UnregisterHandler(name)        
        bucketMtx.Lock()
        defer bucketMtx.Unlock()
        
        bucket := newBucket(handler)
        buckets[name] = bucket
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
        
        var jsonData, textData []byte
        for _, bucket := range buckets {
            if bucket.enabled() {
                switch bucket.formatterType() {
                case JSONFormatter:
                    if jsonData == nil {
                        jsonData = entry.ToJSON()
                    }
                    bucket.entryChan <- jsonData
                case TextFormatter:
                    if textData == nil {
                        textData = entry.ToText()
                    }
                    bucket.entryChan <- textData
                default:
                    bucket.entryChan <- entry
                }
            }
        }
    }    
}
