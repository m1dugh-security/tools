package utils

import (
    "errors"
    "sync"
)

const extensionSize int = 10

type Queue[T any] struct {
    enqueueIndex int
    dequeueIndex int
    values []T
    _arraylen int
    length int
    mut     *sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
    queue := &Queue[T]{
        enqueueIndex: 0,
        dequeueIndex: 0,
        values: make([]T, extensionSize),
        _arraylen: extensionSize,
        length: 0,
        mut: &sync.Mutex{},
    }
    return queue
}

func (q *Queue[T]) _getElements() []T {
    res := make([]T, q.length)
    if q.dequeueIndex < q.enqueueIndex {
        for i := q.dequeueIndex; i < q.enqueueIndex;i++ {
            res[i] = q.values[i]
        }
    } else {
        index := 0
        for i := q.dequeueIndex; i < q._arraylen;i++ {
            res[index] = q.values[i]
            index++
        }

        for i := 0; i < q.enqueueIndex; i++ {
            res[index] = q.values[i]
            index++
        }
    }

    return res
}

func (q *Queue[T]) _shrink() {
    elements := q._getElements()
    q.values = elements
    q._arraylen = len(elements)
    q.length = q._arraylen
    q.dequeueIndex = 0
    q.enqueueIndex = 0
}

func (q *Queue[T]) _flatten() {
    elements := q._getElements()
    copy(q.values, elements)
    q.dequeueIndex = 0
    q.enqueueIndex = len(elements) % q._arraylen
}

func (q *Queue[T]) _extend(deltasize int) {
    freespace := q._arraylen - q.length
    required := deltasize - freespace
    if required <= 0 {
        return
    }

    q._flatten()
    q.values = append(q.values, make([]T, required)...)
    q.enqueueIndex = q._arraylen
    q._arraylen += required
}

func (q *Queue[T]) Enqueue(x T) {
    q.mut.Lock()
    defer q.mut.Unlock()
    if q.length == q._arraylen {
        q._extend(extensionSize)
    }
    q.length++
    q.values[q.enqueueIndex] = x
    q.enqueueIndex = (q.enqueueIndex + 1) % q._arraylen
}

func (q *Queue[T]) Dequeue() (T, error) {
    q.mut.Lock()
    defer q.mut.Unlock()
    var res T
    if q.length == 0 {
        return res, errors.New("Could not dequeue empty queue")
    }
    res = q.values[q.dequeueIndex]
    q.length--
    q.dequeueIndex = (q.dequeueIndex + 1) % q._arraylen
    return res, nil
}

func (q *Queue[T]) Length() int {
    q.mut.Lock()
    defer q.mut.Unlock()
    res := q.length
    return res
}
