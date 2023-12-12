package main

import (
	"reflect"
	"sync"
	"time"
)

type Entry_t map[string]interface{}

type Entry struct {
	data     *Entry_t
	mapping  int
	creation time.Time
	updated  time.Time
}

func StructToMap(input interface{}) *Entry_t {
	result := make(Entry_t)

	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return &result
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return &result
}

func CompareMaps(map1, map2 *Entry_t) bool {
	return reflect.DeepEqual(map1, map2)
}

func BasicEntry(data interface{}) *Entry {
	entryData := StructToMap(data)
	return &Entry{entryData, 0, time.Now(), time.Now()}
}

func (e *Entry) Copy(dest *Entry) *Entry {
	e2 := *e
	if dest != nil {
		e2.data = e.data
		dest.data = &Entry_t{}

		for k, v := range *e.data {
			(*dest.data)[k] = v
		}
	}

	return &e2
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
		if CompareMaps(v.data, entry.data) {
			delete(t.inmem, k)
			break
		}
	}
}

func (u *Entry) getMember(fieldName string) interface{} {
	return (*u.data)[fieldName]
}

func (t *Table) GetWhere(fieldName string, value interface{}) ([]*Entry, int) {
	t.mux.Lock()
	defer t.mux.Unlock()

	res := make([]*Entry, 0)
	amt := 0

	for _, v := range t.inmem {
		field := v.getMember(fieldName)

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
