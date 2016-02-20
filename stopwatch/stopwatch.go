package stopwatch

import (
    "time"
)

type Watch struct {
    start time.Time
    stop time.Time
}

func New() *Watch {
    now := time.Now()
    return &Watch {
        start: now,
        stop: now,
    }
}

func (w *Watch) Stop() {
    w.stop = time.Now()
}

func (w *Watch) Start() {
    w.stop = w.start
}

func (w *Watch) Restart() {
    w.start = time.Now()
    w.stop = w.start
}

func (w *Watch) Duration() time.Duration {
    return w.stop.Sub(w.start)
}

func (w *Watch) Milliseconds() int64 {
    return w.Nanoseconds() / int64(time.Millisecond)
}

func (w *Watch) Microseconds() int64 {
    return w.Nanoseconds() / int64(time.Microsecond)
}

func (w *Watch) Nanoseconds() int64 {
    return w.Duration().Nanoseconds()
}