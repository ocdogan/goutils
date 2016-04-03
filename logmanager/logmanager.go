package logmanager

import (
    "sync/atomic"
)

var isEnabled = int32(0)

// Enable the logging manager globally
func Enable() {
    atomic.StoreInt32(&isEnabled, int32(1))
}

// Disable the logging manager globally
func Disable() {
    atomic.StoreInt32(&isEnabled, int32(0))
}

// Enabled function is used to get if the global logging manager is enabled
func Enabled() bool {
    return atomic.LoadInt32(&isEnabled) == 1
}

// LogError function is used to log the given error by log manager
func LogError(e error) {
    if e != nil && Enabled() {
        err := WrapErrorWithStack(e)
    }
}