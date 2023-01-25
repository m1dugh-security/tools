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

func TestStringSetDiff(t *testing.T) {
    old := NewStringSet("a", "b", "c", "f")
    newSet := NewStringSet("b", "c", "d", "e")

    added, removed := old.Diff(newSet)

    expectedAdded := NewStringSet("d", "e")
    expectedRemoved := NewStringSet("a", "f")

    if !added.Equals(expectedAdded) {
        t.Errorf("Expected %s, but got %s", expectedAdded.UnderlyingArray(), added.UnderlyingArray())
    }
    if !removed.Equals(expectedRemoved) {
        t.Errorf("Expected %s, but got %s", expectedRemoved.UnderlyingArray(), removed.UnderlyingArray())
    }
}

func TestComparableSet(t *testing.T) {
    set := NewComparableSet[testVal]()
    t1 := testVal{ 3, "test1"}
    t2 := testVal{ 2, "test0"}
    t3 := testVal{ 4, "test2"}
    set.AddElement(t1)
    set.AddElement(t2)
    set.AddElement(t3)

    expected := []testVal{ t2, t1, t3 }
    res := set.ToArray()
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("value mismatch: \nreceived: %s\nexpected: %s", res, expected)
    }
}

func TestStringSetSingleThread(t *testing.T) {
    set := NewStringSet()
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

func TestStringSetRemove(t *testing.T) {
    removed := "hello"
    other := "world"
    length := 10
    set := NewStringSet(removed, other)
    for i := 0; i < length - 2; i++ {
        set.AddWord(fmt.Sprintf("test %d", i))
    }


    set.Remove(removed)
    if set.ContainsWord(removed) {
        t.Errorf("Expected set not to contain %s", removed)
    }

    if set.Length() != length - 1 {
        t.Errorf("Expected length to be %d, but got %d", length, set.Length())
    }

    if !set.ContainsWord(other) {
        t.Errorf("Expected set to contain %s", removed)
    }
}

func TestStringSetMultiThread(t *testing.T) {
    set := NewStringSet()
    var wg sync.WaitGroup
    sample := 10000

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(set *StringSet) {
            for i := 0; i < sample; i++ {
                set.AddWord(fmt.Sprintf("test: %d", i))
            }
            wg.Done()
        }(set)
    }
    wg.Wait()

    if set.Length() != sample {
        t.Errorf("Expected the length to be %d but got %d", sample, set.Length())
    }

    kept := 5

    values := set.ToArray()[kept:]
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(set *StringSet) {
            for _, value := range values {
                set.Remove(value)
            }
            wg.Done()
        }(set)
    }
    wg.Wait()

    if set.Length() != kept {
        t.Errorf("Expected the length to be %d but got %d", sample, set.Length())
    }
}

