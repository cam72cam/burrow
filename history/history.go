package history

import (
	"fmt"
	"strings"
)

type Entry struct {
	Name string
	Args []string
}

func (e Entry) String() string {
	return fmt.Sprintf("%s %s", e.Name, strings.Join(e.Args, " "))
}

func NewEntry(input string) *Entry {
	sp := strings.Split(input, " ")
	if len(sp) > 0 && len(strings.TrimSpace(sp[0])) > 0 {
		e := Entry{Name: sp[0]}
		if len(sp) > 1 {
			e.Args = sp[1:]
		}
		return &e
	}
	return nil
}

var entries = make([]Entry, 0)

func Add(e Entry) {
	entries = append(entries, e)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Handles positive numbers as index
// and negative numbers as reverse index
func Get(pos int) *Entry {
	if len(entries) == 0 {
		return nil
	}
	if pos < 0 {
		index := len(entries) + pos
		return &entries[max(index, 0)]
	}
	return &entries[min(pos, len(entries))]
}

func Size() int {
	return len(entries)
}
