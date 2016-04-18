package logmanager

import (
    "fmt"
    "sync/atomic"
)

type ConsoleLogHandler struct {
    disabled uint32
}

func (handler *ConsoleLogHandler) Name() string {
    return "console"
}

func (handler *ConsoleLogHandler) Enabled() bool {
    return atomic.LoadUint32(&handler.disabled) == falseUint32
}

func (handler *ConsoleLogHandler) Enable() {
    atomic.StoreUint32(&handler.disabled, falseUint32)
}

func (handler *ConsoleLogHandler) Disable() {
    atomic.StoreUint32(&handler.disabled, trueUint32)
}

func (handler *ConsoleLogHandler) FormatterType() LogFormatterType {
    return TextFormatter
}

func (handler *ConsoleLogHandler) ProcessText(entry []byte) {
    fmt.Println(string(entry))
}

func (handler *ConsoleLogHandler) ProcessJson(entry []byte) {    
}

func (handler *ConsoleLogHandler) ProcessCustom(entry *LogEntry) {    
}

