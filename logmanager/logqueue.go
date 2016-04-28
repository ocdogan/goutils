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
    "math"
    "sync"
    "sync/atomic"
    "github.com/ocdogan/goutils/utils"
)

const (
    maxQueueLen = int(math.MaxInt16)
)

type logQueueItem struct {
    data interface{}
    next *logQueueItem
}

type logQueue struct {
    sync.Mutex
    cnt int32
    cap int32
    realCap int32
    head *logQueueItem
    tail *logQueueItem
}

func newLogQueue(cap int) *logQueue {
    result := &logQueue{}
    result.setCapacity(cap)
    return result
}

func (q *logQueue) count() int {
    return int(atomic.LoadInt32(&q.cnt))
}

func (q *logQueue) capacity() int {
    return int(atomic.LoadInt32(&q.cap))
}

func (q *logQueue) setCapacity(cap int) {
    if cap < -1 {
        cap = -1
    } else if cap > maxQueueLen {
        cap = maxQueueLen
    }
    q.Lock()
    q.cap = int32(cap)
    q.realCap = q.cap
    q.Unlock()
}

func (q *logQueue) push(data interface{}) {
    if utils.HasValue(data) {
        item := &logQueueItem{
            data: data,
        }

        q.Lock()
        defer q.Unlock()
        
        if q.cnt == q.realCap {
            q.head, q.head.next = q.head.next, nil
            if q.head == nil {
                q.tail = nil
            }
            q.cnt--
        } 
        
        if q.tail == nil {
            q.head = item
            q.tail = item
            return
        }
        q.tail.next = item
        q.tail = item
        q.cnt++
    }
}

func (q *logQueue) pop() (data interface{}) {
    q.Lock()
    defer q.Unlock()
    
    if q.cnt > 0 {
        item := q.head
        q.head = q.head.next
        item.next = nil
        if q.head == nil {
            q.tail = nil
        }
        q.cnt--
        
        data = item.data
        item.data = nil
    }    
    return
}