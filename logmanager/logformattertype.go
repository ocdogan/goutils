package logmanager

type LogFormatterType byte

const (
    TextFormatter LogFormatterType = iota
    JsonFormatter
    CustomFormatter
)