package main

import (
	"reflect"
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

func (e *Entry) Copy(dest *Entry) *Entry {
    e2 := *e
    if dest != nil { dest = &e2 }

    return &e2
}

type UData interface {
//  isEqualTo(d UData) bool
//  getMember(fieldName string) interface{}
}

type offset_t struct {
	Size   uintptr
	Offset uintptr
}

type Offset map[string]offset_t

func CalculateOffsets(iface interface{}) Offset {
	t := reflect.TypeOf(iface)
	offsets := make(Offset)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		offset := uintptr(field.Offset)
		size := uintptr(field.Type.Size())

		offsets[field.Name] = offset_t{Size: size, Offset: offset}
	}

	return offsets
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
