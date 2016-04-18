package uuid

import (
    "crypto/rand"
    "encoding/binary"
    "encoding/hex"
    "net"
    "os"
    "sync"
    "strings"
    "time"
)

const (
    sep = byte('-')
)

type UUID [16]byte

func (uuid UUID) String() string {
	result := make([]byte, 36)

	result[8] = sep
	result[13] = sep
	result[18] = sep
	result[23] = sep

	hex.Encode(result[0:8], uuid[0:4])
	hex.Encode(result[9:13], uuid[4:6])
	hex.Encode(result[14:18], uuid[6:8])
	hex.Encode(result[19:23], uuid[8:10])
	hex.Encode(result[24:], uuid[10:])

	return strings.ToUpper(string(result))
}

type seqID struct {
    sync.Mutex
    id uint64
}

var (
    now uint64
    euid []byte
    haddr []byte
    id = &seqID{}
)

func NewUUID() (*UUID, error) {
    uuid := new(UUID)
    _, err := rand.Read(uuid[:])
    if err != nil {
        return nil, err
    }
    
    uuid.xor(euid, 0)    
    uuid.xor(getDate(), 8)
    uuid.xor(haddr, 4)
    
    return uuid, nil
}

func (uuid *UUID) xor(bytes []byte, shift int) {
    for i, b := range bytes {
        pos := i + shift
        if pos >= 16 {
            break
        }
        uuid[pos] ^= b 
    }
}

func getDate() []byte {    
    id.Lock()
    id.id++
    tick := id.id
    id.Unlock()
    
    result := make([]byte, 8) 
    binary.BigEndian.PutUint64(result, now + tick)

    return result
}

func init() {
    defer recover()
    
    now = uint64(time.Now().Unix())
    
    euid = make([]byte, 8) 
    binary.BigEndian.PutUint64(euid, uint64(os.Getpid()))    
    
    nowArr := make([]byte, 8) 
    binary.BigEndian.PutUint64(nowArr, now)
    
    for i, b := range euid {
        euid[i] = b ^ nowArr[i]
    }
    
    initHAddr()    
}

func initHAddr() {
    if haddr != nil {
        return
    }
    
    infs, err := net.Interfaces()
    if err != nil {
        return
    }

    var lenAdd int    
    for _, inf := range infs {
        lenH := len(inf.HardwareAddr)
        if lenH > 6 {
            if haddr == nil {
                lenAdd = lenH
                haddr = make([]byte, lenH)
                copy(haddr, inf.HardwareAddr)
            } else {
                l := lenH
                if lenAdd < l {
                    l = lenAdd
                }
                
                for i, b := range inf.HardwareAddr {
                    if i >= l {
                        break
                    }
                    haddr[i] ^= b 
                }
            }
        }
    }
}