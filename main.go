package main

import (
    "fmt"
    "time"
    "sync"
)

type Entry struct {
    data       UData
    creation   time.Time
    updated    time.Time
    member_off []byte
}

type UData interface {
    getMember(s string) interface{}
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
    Name     string
    Age      int
}

func (e Employee) getMember(s string) interface{} {
    switch s {
    case "Name":
        return e.Name;
    case "Age":
        return e.Age;
    }

    return nil;
}


func main() {
    table := Table{ inmem: make(map[int64]*Entry) }

    table.Add(&(Entry{ Employee { "albert einstein", 144 }, time.Now(), time.Now(), make([]byte, 0)}))
    table.Add(&(Entry{ Employee { "albert einstein", 10 }, time.Now(), time.Now(), make([]byte, 0)}))

    query, amt := table.GetByField("Name", "albert einstein")

    fmt.Println("found", amt, "entries")
    
    for _, v := range query {
        fmt.Println(*v)
    }
}
