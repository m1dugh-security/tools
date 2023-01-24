package utils

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type testVal struct {
    val int
    name string
}

func (v testVal) String() string {
    return fmt.Sprintf("%d: %s", v.val, v.name)
}

func (t testVal) Compare(other interface{}) int {
    t2, ok := other.(testVal)
    if ok {
        return t.val - t2.val
    } else {
        return t.val
    }
}

func TestComparableSet(t *testing.T) {
    var set *ComparableSet[testVal] = new(ComparableSet[testVal])
    set.AddElement(testVal{
        3,
        "test1",
    })

    set.AddElement(testVal{
        2,
        "test0",
    })

    set.AddElement(testVal{
        4,
        "test2",
    })
    expected := []testVal{
        testVal{
            2,
            "test0",
        },
        testVal{
            3,
            "test1",
        },
        testVal{
            4,
            "test2",
        },
    }
    res := set.ToArray()
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("value mismatch: \nreceived: %s\nexpected: %s", res, expected)
    }
}

func TestStringSetSingleThread(t *testing.T) {
    set := NewStringSet(nil)
    word := "hello"
    if set.AddWord(word) != true {
        t.Errorf("Expected to be able to add word")
    }

    for i := 0; i < 1000; i++ {
        if set.AddWord(word) != false {
            t.Errorf("Expected to find duplicate")
        }
    }
}

func TestStringSetMultiThead(t *testing.T) {
    set := NewStringSet(nil)
    var wg sync.WaitGroup
    sample := 10000

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(set *StringSet) {
            for i := 0; i < sample; i++ {
                set.AddWord(fmt.Sprintf("%d", i))
            }
            wg.Done()
        }(set)
    }
    wg.Wait()

    if set.Length() != sample {
        t.Errorf("Expected the length to be %d but got %d", sample, set.Length())
    }
}

