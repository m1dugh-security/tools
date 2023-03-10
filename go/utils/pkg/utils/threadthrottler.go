package utils

import (
    "sync"
)

type ThreadThrottler struct {
    MaxThreads uint
    threads uint
    mut *sync.Mutex
    wg *sync.WaitGroup
}

func NewThreadThrottler(maxThreads uint) *ThreadThrottler {
    return &ThreadThrottler{
        maxThreads,
        0,
        &sync.Mutex{},
        &sync.WaitGroup{},
    }
}

func (t *ThreadThrottler) RequestThread() {
    t.mut.Lock()
    defer t.mut.Unlock()
    t.threads++
    if t.threads <= t.MaxThreads {
        t.wg.Add(1)
        return
    }
    t.mut.Unlock()
    for true {
        t.mut.Lock()
        if t.threads <= t.MaxThreads {
            t.wg.Add(1)
            break
        }
        t.mut.Unlock()
    }
}

func (t *ThreadThrottler) Threads() uint {
    t.mut.Lock()
    defer t.mut.Unlock()
    return t.threads
}

func (t *ThreadThrottler) Done() {

    t.mut.Lock()
    if t.threads > 0 {
        t.threads--
    }
    t.mut.Unlock()
    t.wg.Done()
}

func (t *ThreadThrottler) Wait() {
    t.wg.Wait()
}

