package utils

import (
    "sync"
    "time"
)

type RequestThrottler struct {
    MaxRequests int
    _requests int
    _lastFlush int64
    mut *sync.Mutex
}

func NewRequestThrottler(maxRequests int) *RequestThrottler {
    res := &RequestThrottler{
        MaxRequests: maxRequests,
        _requests: 0,
        _lastFlush: time.Now().UnixMicro(),
        mut: &sync.Mutex{},
    }
    return res
}

func (r *RequestThrottler) AskRequest() {
    if r.MaxRequests < 0 {
        return
    }
    r.mut.Lock()
    defer r.mut.Unlock()

    timeStampMicro := time.Now().UnixMicro()
    delta := timeStampMicro - r._lastFlush
    if delta > 1000000 {
        r._requests = 1
        r._lastFlush = timeStampMicro
    } else {
        if r._requests < r.MaxRequests {
            r._requests++
        } else {
            for timeStampMicro - r._lastFlush < 1000000 {
                time.Sleep(time.Microsecond)
                timeStampMicro = time.Now().UnixMicro()
            }

            r._requests = 1
            r._lastFlush = timeStampMicro
        }
    }
}
