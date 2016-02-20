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