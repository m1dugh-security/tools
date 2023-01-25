package utils

import (
	"sync"
)

type StringSet struct {
    values  []string
    mut     *sync.Mutex
}

func NewStringSet(values... string) *StringSet {
    
    var res *StringSet = &StringSet{
        values: make([]string, 0, len(values)),
        mut: &sync.Mutex{},
    }

    for _, v := range values {
        res.addWord(v)
    }

    return res
}

func (set *StringSet) length() int {
    if set == nil {
        return 0
    }
    return len(set.values)
}

func (set *StringSet) Length() int {
    set.mut.Lock()
    defer set.mut.Unlock()
    return set.length()
}

func (set *StringSet) remove(value string) bool {
    pos, found := set._binsearch(value)
    if !found {
        return false
    }

    middle := set.length() / 2
    if pos > middle {
        for i := pos; i < set.length() - 1; i++ {
            set.values[i] = set.values[i + 1]
        }
        set.values = set.values[:set.length() - 1]
    } else {
        for i := pos - 1; i > 0 ; i-- {
            set.values[i + 1] = set.values[i]
        }
        set.values = set.values[1:]
    }

    return true
}

func (set *StringSet) Remove(value string) bool {
    set.mut.Lock()
    defer set.mut.Unlock()
    return set.remove(value)
}

/// Returns the underlying array.
/// Caution: Since it does not return a copy, modifying the underlying array
/// can break the StringSet
func (set *StringSet) UnderlyingArray() []string {
    return set.values
}

func (set *StringSet) _binsearch(value string) (int, bool) {

    start := 0
    end := set.length()

    for start < end {
        middle := start + (end - start) / 2
        s := set.values[middle]
        res := CompareStrings(value, s)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *StringSet) _insertAt(value string, pos int) {
    set.values = append(set.values, value)
    for i := set.length() - 1; i > pos; i-- {
        set.values[i] = set.values[i - 1]
    }

    set.values[pos] = value
}
func (set *StringSet) addWord(value string) bool {
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *StringSet) AddWord(value string) bool {
    set.mut.Lock()
    defer set.mut.Unlock()
    return set.addWord(value)
}

func (set *StringSet) containsWord(value string) bool {
    _, found := set._binsearch(value)
    return found
}

func (set *StringSet) ContainsWord(value string) bool {
    set.mut.Lock()
    defer set.mut.Unlock()
    return set.containsWord(value)
}

func (set *StringSet) AddAll(other *StringSet) int {
    set.mut.Lock()
    defer set.mut.Unlock()
    var count int = 0
    for _, v := range other.values {
        if set.addWord(v) {
            count++
        }
    }

    return count
}

func (set *StringSet) ToArray() []string {
    set.mut.Lock()
    defer set.mut.Unlock()
    res := make([]string, set.length())
    copy(res, set.values)
    return res
}

/// calculates the diff between two string sets
/// returns a Set containg the strings added in the newSet, 
/// and a set containing the strings removed in the newSet
func (old *StringSet) Diff(newSet *StringSet) (*StringSet, *StringSet) {
    if old == nil {
        return newSet, NewStringSet()
    } else if newSet == nil {
        return NewStringSet(), old
    }
    old.mut.Lock()
    defer old.mut.Unlock()
    newSet.mut.Lock()
    defer newSet.mut.Unlock()
    added := NewStringSet()
    removed := NewStringSet()

    oldi := 0
    newi := 0
    for oldi < old.length() && newi < newSet.length() {
        s1 := old.values[oldi]
        s2 := newSet.values[newi]

        res := CompareStrings(s1, s2)
        if res < 0 {
            oldi++
            removed.addWord(s1)
        } else if res > 0 {
            newi++
            added.addWord(s2)
        } else {
            oldi++
            newi++
        }
    }

    for ; newi < newSet.length(); newi++ {
        s2 := newSet.values[newi]
        added.addWord(s2)
    }

    for ; oldi < old.length(); oldi++ {
        s1 := old.values[oldi]
        removed.addWord(s1)
    }

    return added, removed
}

func (s *StringSet) Equals(other *StringSet) bool {
    s.mut.Lock()
    s.mut.Unlock()
    if s.length() != other.length() {
        return false
    }

    for i := 0; i < s.length(); i++ {
        if s.values[i] != other.values[i] {
            return false
        }
    }

    return true
}

type Comparable interface {
    Compare(value interface{}) int
}

type ComparableSet[T Comparable] struct {
    values []T
    mut *sync.Mutex
}

func NewComparableSet[T Comparable](values... T) *ComparableSet[T] {

    res := &ComparableSet[T]{
        values: make([]T, 0, len(values)),
        mut: &sync.Mutex{},
    }

    for _, v := range values {
        res.addElement(v)
    }

    return res
}

func (set *ComparableSet[T]) _binsearch(value T) (int, bool) {

    start := 0
    end := set.length()

    for start < end {
        middle := start + (end - start) / 2
        elem := set.values[middle]
        res := value.Compare(elem)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *ComparableSet[T]) _insertAt(value T, pos int) {
    set.values = append(set.values, value)
    for i := set.length()- 1; i > pos; i-- {
        set.values[i] = set.values[i - 1]
    }

    set.values[pos] = value
}

func (set *ComparableSet[T]) addElement(value T) bool {
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *ComparableSet[T]) AddElement(value T) bool {
    if set == nil {
        return false
    }
    set.mut.Lock()
    defer set.mut.Unlock()

    return set.addElement(value)
}

func (set *ComparableSet[T]) contains(value T) bool {
    _, found := set._binsearch(value)
    return found
}

func (set *ComparableSet[T]) Contains(value T) bool {
    if set == nil {
        return false
    }
    set.mut.Lock()
    defer set.mut.Unlock()

    return set.contains(value)
}

func (set *ComparableSet[T]) AddAll(other *ComparableSet[T]) int {
    set.mut.Lock()
    defer set.mut.Unlock()
    var count int = 0
    for _, v := range other.values {
        if set.addElement(v) {
            count++
        }
    }

    return count
}

func (set *ComparableSet[T]) length() int {
    if set == nil {
        return 0
    }
    return len(set.values)
}

func (set *ComparableSet[T]) Length() int {
    if set == nil {
        return 0
    }
    set.mut.Lock()
    defer set.mut.Unlock()
    return len(set.values)
}

/// Returns the underlying array.
/// Caution: Since it does not return a copy, modifying the underlying array
/// can break the ComparableSet
func (set *ComparableSet[T]) UnderlyingArray() []T {
    return set.values
}

func (set *ComparableSet[T]) ToArray() []T {
    set.mut.Lock()
    defer set.mut.Unlock()
    res := make([]T, set.length())
    copy(res, set.values)
    return res
}

func (old *ComparableSet[T]) Diff(newSet *ComparableSet[T]) (*ComparableSet[T], *ComparableSet[T]) {

    added := NewComparableSet[T]()
    removed := NewComparableSet[T]()

    oldi := 0
    newi := 0
    for oldi < old.length() && newi < newSet.length() {
        s1 := old.values[oldi]
        s2 := newSet.values[newi]

        res := s1.Compare(s2)
        if res < 0 {
            oldi++
            removed.addElement(s1)
        } else if res > 0 {
            newi++
            added.addElement(s2)
        } else {
            oldi++
            newi++
        }
    }

    for ; newi < newSet.length(); newi++ {
        s2 := newSet.UnderlyingArray()[newi]
        added.addElement(s2)
    }

    for ; oldi < old.length(); oldi++ {
        s1 := old.UnderlyingArray()[oldi]
        removed.addElement(s1)
    }

    return added, removed
}

func (s *ComparableSet[T]) Equals(other *ComparableSet[T]) bool {
    if s.length() != other.length() {
        return false
    }

    for i := 0; i < s.length(); i++ {
        v1 := s.UnderlyingArray()[i]
        v2 := other.UnderlyingArray()[i] 
        if v1.Compare(v2) != 0 {
            return false
        }
    }

    return true
}

func (set *ComparableSet[T]) remove(value T) bool {
    pos, found := set._binsearch(value)
    if !found {
        return false
    }

    middle := set.length() / 2
    if pos > middle {
        for i := pos; i < set.length() - 1; i++ {
            set.values[i] = set.values[i + 1]
        }
        set.values = set.values[:set.length() - 1]
    } else {
        for i := pos - 1; i > 0 ; i-- {
            set.values[i + 1] = set.values[i]
        }
        set.values = set.values[1:]
    }

    return true
}

func (set *ComparableSet[T]) Remove(value T) bool {
    set.mut.Lock()
    defer set.mut.Unlock()
    return set.remove(value)
}
