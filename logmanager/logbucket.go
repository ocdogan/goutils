package logmanager

import (
    "bytes"
    "container/list"
    "encoding/json"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
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
    entryChan chan *LogEntry
    queue *list.List
    handler LogHandler
}

func newBucket(handler LogHandler) *logBucket {
    return &logBucket{
        handler: handler,
        queue: list.New(),
        done: make(chan bool),
        entryChan: make(chan *LogEntry),
        completed: falseUint32,
    }
}

func (bucket *logBucket) enabled() bool {
    return atomic.LoadUint32(&bucket.completed) == falseUint32 && 
        bucket.handler != nil && 
        bucket.handler.Enabled()
}

func (bucket *logBucket) close() {
    bucket.done <- true
    close(bucket.done)
    close(bucket.entryChan)
}
    
func (bucket *logBucket) pop() *LogEntry {
    bucket.Lock()
    defer bucket.Unlock()
    
    if bucket.queue.Len() > 0 {
        elm := bucket.queue.Back()
        if elm != nil {
            bucket.queue.Remove(elm)
            return elm.Value.(*LogEntry)
        }
    }
    return nil
}

func (bucket *logBucket) push(entry *LogEntry) {
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

    bucket.queue.PushBack(entry)
}

func formatAsText(entry *LogEntry) []byte {
    if entry == nil {
        return nil
    }
    
    buffer := &bytes.Buffer{}

    writeToBuffer(buffer, "id", entry.id)
    writeToBuffer(buffer, "time", entry.time)
    writeToBuffer(buffer, "duration", entry.duration)
    writeToBuffer(buffer, "logType", entry.logType.String())
    writeToBuffer(buffer, "message", entry.message)
    writeToBuffer(buffer, "stack", entry.stack)
    
    if entry.args == nil {
        writeToBuffer(buffer, "args=\"\"", nil)    
    } else {
        for k, v := range entry.args {
            writeToBuffer(buffer, k, v)
        }
    }
    
    return buffer.Bytes()
}

func isTextOrNumber(value string) bool {
	for _, c := range value {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
            c == '-' || c == '+' || 
            c == '.' || c == ',') {
			return true
		}
	}
	return false
}

func writeToBuffer(buffer *bytes.Buffer, key string, value interface{}) {
    buffer.WriteString(key)
    
    switch value.(type) {
    case nil:
        break
    case string:
        if !isTextOrNumber(value.(string)) {
            fmt.Fprintf(buffer, "=%q ", value)
        } else {
            buffer.WriteString("=\"")
            buffer.WriteString(value.(string))
            buffer.WriteString("\" ")
        }
    default:
        buffer.WriteString("=\"")
        fmt.Fprint(buffer, value)         
        buffer.WriteString("\" ")
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
            entry := bucket.pop() 
            if entry == nil {
                continue
            }
            handler := bucket.handler
            if handler.Enabled() {
                var jsonEntry, textEntry []byte
                switch handler.FormatterType() {
                case JsonFormatter:
                    if jsonEntry == nil {
                        b, e := json.Marshal(struct{
                            ID string `json:"id"`
                            Time time.Time `json:"time"`
                            Duration time.Duration `json:"duration"`
                            LogType string `json:"logType"`
                            Message string `json:"message"`
                            Stack string `json:"stack"`
                            Args map[string]interface{} `json:"args"`
                        }{
                            ID: entry.id,
                            Time: entry.time,
                            Duration: entry.duration,
                            LogType: entry.logType.String(),
                            Message: entry.message,
                            Stack: entry.stack,
                            Args: entry.args,
                        })
                        if e != nil {
                            continue
                        }
                        jsonEntry = b
                    }
                    handler.ProcessJson(jsonEntry)
                case TextFormatter:
                    if textEntry == nil {                        
                        textEntry = formatAsText(entry)
                    }
                    handler.ProcessText(textEntry)
                case CustomFormatter:
                    handler.ProcessCustom(entry)
                }
            }
        }
    }
}
