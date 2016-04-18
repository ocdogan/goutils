package logmanager

// LogHandler is the interface which handles log writes to different destinations
type LogHandler interface {
    Name() string
    Enable()
    Disable()
    Enabled() bool
    FormatterType() LogFormatterType
    ProcessJson(entry []byte)
    ProcessText(entry []byte)
    ProcessCustom(entry *LogEntry)
}