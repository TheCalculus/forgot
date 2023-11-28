package main

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	data     UData
	mapping  int
	creation time.Time
	updated  time.Time
}

func BasicEntry(data UData) *Entry {
	return &Entry{data, 0, time.Now(), time.Now()}
}

type UData interface {
	getMember(s string) interface{}
	isEqualTo(d UData) bool
}

type Table struct {
	mux    sync.Mutex
	inmem  map[int]*Entry
	active int
}

func (t *Table) Add(entry *Entry) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.inmem[t.active] = entry
	entry.mapping = t.active

	t.active++
}

func (t *Table) Remove(entry *Entry, mult bool) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !mult {
		delete(t.inmem, entry.mapping)
		return
	}

	for k, v := range t.inmem {
		if v.data.isEqualTo(entry.data) {
			delete(t.inmem, k)
		}
	}
}

func (t *Table) GetWhere(fieldName string, value interface{}) ([]*Entry, int) {
	t.mux.Lock()
	defer t.mux.Unlock()

	res := make([]*Entry, 0)
	amt := 0

	for _, v := range t.inmem {
		field := v.data.getMember(fieldName)

		if field == nil {
			return nil, -1
		}

		if field == value {
			res = append(res, v)
			amt++
		}
	}

	return res, amt
}

// user-defined data here
type Employee struct {
	Name string
	Age  int
}

func (e Employee) getMember(s string) interface{} {
	switch s {
	case "Name":
		return e.Name
	case "Age":
		return e.Age
	}

	return nil
}

func (e Employee) isEqualTo(d UData) bool {
	if other, ok := d.(Employee); ok {
		return e.Name == other.Name && e.Age == other.Age
	}
	return false
}

func newEmployee(name string, age int) Employee {
	return Employee{name, age}
}

func main() {
	table := Table{inmem: make(map[int]*Entry)}

	albert := BasicEntry(newEmployee("albert einstein", 144))

	table.Add(albert)
	table.Add(BasicEntry(newEmployee("albert einstein", 144)))
	table.Add(BasicEntry(newEmployee("albert einstein", 10)))

	table.Remove(albert, true)

	query, _ := table.GetWhere("Name", "albert einstein")

	for _, v := range query {
		fmt.Println(*v)
	}
}
