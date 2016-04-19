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
    "container/list"
    "sync"
    "sync/atomic"
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
    completed uint32
    entryChan chan interface{}
    queue *list.List
    handler LogHandler
}

func newBucket(handler LogHandler) *logBucket {
    return &logBucket{
        handler: handler,
        queue: list.New(),
        done: make(chan bool),
        entryChan: make(chan interface{}),
        completed: falseUint32,
    }
}

func (bucket *logBucket) enabled() bool {
    return atomic.LoadUint32(&bucket.completed) == falseUint32 && 
        bucket.handler != nil && 
        bucket.handler.Enabled()
}

func (bucket *logBucket) formatterType() LogFormatterType {
    if bucket.handler != nil { 
        return bucket.handler.FormatterType()
    }
    return CustomFormatter
}

func (bucket *logBucket) close() {
    bucket.done <- true
    close(bucket.done)
    close(bucket.entryChan)
}
    
func (bucket *logBucket) pop() interface{} {
    bucket.Lock()
    defer bucket.Unlock()
    
    if bucket.queue.Len() > 0 {
        elm := bucket.queue.Back()
        if elm != nil {
            bucket.queue.Remove(elm)
            return elm.Value
        }
    }
    return nil
}

func (bucket *logBucket) push(data interface{}) {
    bucket.Lock()
    defer bucket.Unlock()
    
    l := bucket.queue.Len()
    cap := int(BucketCapacity())

    for l > 0 && l >= cap {
        elm := bucket.queue.Back()
        if elm == nil {
            break
        }
        bucket.queue.Remove(elm)
        l = bucket.queue.Len()
    }

    switch data.(type) {
    case []byte:
        bucket.queue.PushBack(data)
    case *LogEntry:
        bucket.queue.PushBack(data)
    case LogEntry:
        bucket.queue.PushBack(&data)
    case string:
        s := data.(string)
        if s != "" {
            bucket.queue.PushBack([]byte(s))
        }
    }
}

func (bucket *logBucket) process() {
    for bucket.enabled() {
        select {
        case <-bucket.done:
            atomic.StoreUint32(&bucket.completed, trueUint32)
            return
        case entry := <-bucket.entryChan:
            bucket.push(entry)
        default:
            e := bucket.pop()
            if e == nil {
                continue
            }
            handler := bucket.handler
            if handler.Enabled() {
                switch handler.FormatterType() {
                case JsonFormatter:
                    b, ok := e.([]byte)
                    if ok && len(b) > 0 {
                        handler.ProcessJson(b)
                    }
                case TextFormatter:
                    b, ok := e.([]byte)
                    if ok && len(b) > 0 {
                        handler.ProcessText(b)
                    }
                case CustomFormatter:
                    entry, ok := e.(*LogEntry)
                    if ok && entry != nil {
                        handler.ProcessCustom(entry)
                    }
                }
            }
        }
    }
}
