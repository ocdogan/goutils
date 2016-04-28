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
    // "container/list"
    "sync"
    "sync/atomic"
    "github.com/ocdogan/goutils/utils"
)

const (
    minBucketCap = uint32(8)
    maxBucketCap = uint32(16*1024)
)

var (
    bucketCap = uint32(1024)
)

// BucketCapacity returns the system wide LogEntry buffering capacity for any log handler
func BucketCapacity() uint32 {
    return atomic.LoadUint32(&bucketCap)
}

// SetBucketCapacity sets the system wide LogEntry buffering capacity for any log handler
func SetBucketCapacity(cap uint32) {
    if cap < minBucketCap {
        cap = minBucketCap
    } else if cap > maxBucketCap {
        cap = maxBucketCap
    }
    
    atomic.StoreUint32(&bucketCap, cap)
}

type logBucket struct {
    sync.Mutex    
    done chan bool
    inProc uint32
    completed uint32
    readyChan chan bool
    queueChan chan interface{}
    queue *logQueue
    handler LogHandler
}

func newBucket(handler LogHandler) *logBucket {
    return &logBucket{
        handler: handler,
        queue: newLogQueue(handler.QueueLen()),
        done: make(chan bool),
        readyChan: make(chan bool),
        queueChan: make(chan interface{}),
        completed: falseUint32,
    }
}

func (bucket *logBucket) enabled() bool {
    return atomic.LoadUint32(&bucket.completed) == falseUint32 && 
        bucket.handler != nil && 
        bucket.handler.Enabled()
}

func (bucket *logBucket) processing() bool {
    return atomic.LoadUint32(&bucket.inProc) != falseUint32
}

func (bucket *logBucket) level() LogLevel {
    if bucket.handler != nil { 
        l := bucket.handler.Level()
        if l != LogLevel(0) {
            return l
        }
    }
    return AllLogLevels
}

func (bucket *logBucket) format() LogFormat {
    if bucket.handler != nil { 
        return bucket.handler.Format()
    }
    return CustomFormat
}

func (bucket *logBucket) queueLen() int {
    if bucket.handler != nil { 
        ql := bucket.handler.QueueLen()
        if ql < 0 {
            return int(BucketCapacity())
        }
        if ql < int(minBucketCap) {
            ql = int(minBucketCap)
        } else if ql > int(maxBucketCap) {
            ql = int(maxBucketCap)
        }
        return ql
    }
    return 0
}

func (bucket *logBucket) close() {
    bucket.done <- true
    close(bucket.done)
    close(bucket.readyChan)
    close(bucket.queueChan)
}
    
func (bucket *logBucket) push(data interface{}) {
    switch data.(type) {
    case []byte:
        bucket.queue.push(data)
    case *LogEntry:
        bucket.queue.push(data)
    case LogEntry:
        bucket.queue.push(&data)
    case string:
        s := data.(string)
        if s != "" {
            bucket.queue.push([]byte(s))
        }
    }
}

func (bucket *logBucket) process() {
    for bucket.enabled() {
        select {
        case <-bucket.done:
            atomic.StoreUint32(&bucket.completed, trueUint32)
            return
        case ready, ok := <-bucket.readyChan:
            if !ok {
                return
            }
            if ready {
                e := bucket.queue.pop()
                if !utils.HasValue(e) {
                    continue
                }
                handler := bucket.handler
                if handler.Enabled() {
                    switch handler.Format() {
                    case TextFormat:
                        fallthrough
                    case JSONFormat:
                        b, ok := e.([]byte)
                        if ok && len(b) > 0 {
                            handler.Process(b)
                        }
                    case CustomFormat:
                        entry, ok := e.(*LogEntry)
                        if ok && entry != nil {
                            handler.Process(entry)
                        }
                    }
                }
            }
        }
    }
}

func (bucket *logBucket) handle() {
    if !atomic.CompareAndSwapUint32(&bucket.inProc, falseUint32, trueUint32) {
        return
    }
    go bucket.process()
    
    for bucket.enabled() {
        select {
        case <-bucket.done:
            atomic.StoreUint32(&bucket.completed, trueUint32)
            return
        case entry, ok := <-bucket.queueChan:
            if !ok {
                return
            }
            bucket.push(entry)
            bucket.readyChan <- true
        }
    }
}
