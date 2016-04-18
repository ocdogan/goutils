package logmanager

type LogType byte

const (
    InfoLog LogType = iota
    WarningLog
    ErrorLog
    FatalLog
)

func (lt LogType) String() string {
    switch lt {
    case InfoLog:
        return "info"
    case WarningLog:
        return "warning"
    case ErrorLog:
        return "error"
    case FatalLog:
        return "fatal"
    }
    return ""
}