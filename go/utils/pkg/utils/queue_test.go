package utils

import (
	"fmt"
	"sync"
	"testing"
)

const sample int = 10000

func TestQueueSingleThread(t *testing.T) {
    q := NewQueue[string]()
    var testData []string
    for i := 0; i < sample; i++ {
        testData = append(testData, fmt.Sprintf("test %d", i))
    }

    for _, v := range testData {
        q.Enqueue(v)
    }

    if q.Length() != sample {
        t.Errorf("Expected length to be %d", sample)
    }
    for i, v := range testData {
        data, err := q.Dequeue()
        if err != nil {
            t.Errorf("Unexpected error while dequeuing")
        }

        if data != v {
            t.Errorf("Wrong data dequeued at index %d: got: %s\nexpected: %s", i, data, v)
        }
    }

    if q.Length() != 0 {
        t.Errorf("Expected length to be %d", 0)
    }
}

func TestQueueMultiThreaded(t *testing.T) {

    q := NewQueue[string]()
    var wg sync.WaitGroup
    var testData []string
    for i := 0; i < sample; i++ {
        testData = append(testData, fmt.Sprintf("test %d", i))
    }

    for _, v := range testData {
        wg.Add(1)
        go func() {
            defer wg.Done()
            q.Enqueue(v)
        }()
    }

    wg.Wait()

    if q.Length() != sample {
        t.Errorf("Expected length to be %d", sample)
    }

    for i := 0; i < sample; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            data, err := q.Dequeue()
            if err != nil {
                t.Errorf("Unexpected error while dequeuing")
            }
            j := 0
            for ; j < sample && data != testData[j]; j++ {}
            if j == sample {
                t.Errorf("Could not find element")
            }
        }()
    }

    wg.Wait()

    if q.Length() != 0 {
        t.Errorf("Expected length to be %d", 0)
    }
}
