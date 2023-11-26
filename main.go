package main

import (
    "fmt"
    "time"
    "sync"
    "reflect"
)

type Entry struct {
    data     interface{}
    creation time.Time
    updated  time.Time
}

type Table struct {
    mux      sync.Mutex
    inmem    map[int64]*Entry
    active   int64
}
    
func (t *Table) Add(entry *Entry) {
    t.mux.Lock()
    defer t.mux.Unlock()
    t.inmem[t.active] = entry
    t.active++
}

func (t *Table) GetByField(fieldName string, value interface{}) ([]*Entry, int) {
    res := make([]*Entry, 0);
    amt := 0

    for _, v := range t.inmem {
        field, exists := getField(v.data, fieldName)

        if !exists {
            return nil, -1
        }

        if field == value {
            res = append(res, v)
            amt++
        }
    }

    return res, amt
}

func getField(obj interface{}, fieldName string) (interface{}, bool) {
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Struct {
		return nil, false
	}

	field := val.FieldByName(fieldName)

	if !field.IsValid() {
		return nil, false
	}

	return field.Interface(), true
}

// user-defined data here
type Employee struct {
    Name     string
    Age      int
}

func main() {
    table := Table{ inmem: make(map[int64]*Entry) }

    table.Add(&(Entry{ Employee { "albert einstein", 144 }, time.Now(), time.Now() }))
    table.Add(&(Entry{ Employee { "albert einstein", 10 }, time.Now(), time.Now() }))

    query, amt := table.GetByField("Name", "albert einstein")

    fmt.Println("found", amt, "entries")
    
    for _, v := range query {
        fmt.Println(*v)
    }
}
