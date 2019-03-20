package er

import (
	"reflect"
	"sort"
	"unsafe"
)

type Query []interface{}

type QueryResult struct {
	cur, max int
	base     unsafe.Pointer
	match    []cmpFn
	s        EntityMeta
	err      error
}

type EntityMeta struct {
	entityWidth uintptr
	cols        []attributeMeta
}

func (s EntityMeta) EvalQuery(q Query, rows interface{}) *QueryResult {
	data, n := s.unsafeTableData(rows)
	key, match, err := s.scanQuery(q)
	if err != nil {
		return &QueryResult{err: err}
	}
	cur := s.followIndex(key, data, 0, n, true)
	max := s.followIndex(key, data, cur, n, false)
	return &QueryResult{
		cur:   cur - 1,
		max:   max,
		base:  data,
		match: match,
		s:     s,
	}
}

func (r *QueryResult) Next() bool {
	if r.err != nil {
		return false
	}
	for r.cur++; r.cur < r.max; r.cur++ {
		if r.currentMatches() {
			return true
		}
	}
	return false
}

func (r *QueryResult) This() int {
	return r.cur
}

func (r *QueryResult) Err() error {
	return r.err
}

type attributeMeta struct {
	offset     uintptr
	key        bool
	columnType func(uintptr, interface{}) (func(unsafe.Pointer) int, error)
}

type cmpFn struct {
	f func(unsafe.Pointer) int
}

func compareInt(offset uintptr, value interface{}) (func(unsafe.Pointer) int, error) {
	x, ok := value.(int)
	if !ok {
		return nil, ErrInvalidAttribute
	}
	return func(p unsafe.Pointer) int {
		y := *(*int)(unsafe.Pointer(uintptr(p) + offset))
		switch {
		case x > y:
			return -1
		case x < y:
			return 1
		}
		return 0
	}, nil
}

func compareStr(offset uintptr, value interface{}) (func(unsafe.Pointer) int, error) {
	x, ok := value.(string)
	if !ok {
		return nil, ErrInvalidAttribute
	}
	return func(p unsafe.Pointer) int {
		y := *(*string)(unsafe.Pointer(uintptr(p) + offset))
		switch {
		case x > y:
			return -1
		case x < y:
			return 1
		}
		return 0
	}, nil
}

func compareFlt(offset uintptr, value interface{}) (func(unsafe.Pointer) int, error) {
	x, ok := value.(float64)
	if !ok {
		return nil, ErrInvalidAttribute
	}
	return func(p unsafe.Pointer) int {
		y := *(*float64)(unsafe.Pointer(uintptr(p) + offset))
		switch {
		case x > y:
			return -1
		case x < y:
			return 1
		}
		return 0
	}, nil
}

func SpecFor(t reflect.Type) EntityMeta {
	var cols []attributeMeta
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		_, key := f.Tag.Lookup("key")
		var cmp func(uintptr, interface{}) (func(unsafe.Pointer) int, error)
		switch f.Type.Kind() {
		case reflect.Int:
			cmp = compareInt
		case reflect.String:
			cmp = compareStr
		case reflect.Float64:
			cmp = compareFlt
		}
		cols = append(cols, attributeMeta{
			offset:     f.Offset,
			key:        key,
			columnType: cmp,
		})
	}
	return EntityMeta{
		entityWidth: t.Size(),
		cols:        cols,
	}
}

func (s EntityMeta) unsafeTableData(rows interface{}) (unsafe.Pointer, int) {
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(reflect.ValueOf(&rows).Elem().InterfaceData()[1]))
	return unsafe.Pointer(sh.Data), sh.Len
}

func (s EntityMeta) scanQuery(q Query) (key, match []cmpFn, err error) {
	match = make([]cmpFn, 0, len(s.cols))
	inKey := true
	for attr, value := range q {
		col := s.cols[attr]
		if attr == len(q)-1 || inKey && (value == nil || !col.key) {
			inKey = false
			key = match
			match = match[attr:attr]
		}
		if value == nil {
			continue
		}
		cmp, err := col.columnType(col.offset, value)
		if err != nil {
			return nil, nil, err
		}
		match = append(match, cmpFn{cmp})
	}
	return key, match, nil
}

func (s EntityMeta) getRecord(data unsafe.Pointer, idx int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(data) + uintptr(idx)*s.entityWidth)
}

func (s EntityMeta) followIndex(key []cmpFn, data unsafe.Pointer, min, max int, deft bool) int {
	return sort.Search(max-min, func(idx int) bool {
		record := s.getRecord(data, idx+min)
		for _, cmp := range key {
			comp := cmp.f(record)
			switch {
			case comp < 0:
				return false
			case comp > 0:
				return true
			}
		}
		return deft
	}) + min
}

func (r *QueryResult) currentMatches() bool {
	record := r.s.getRecord(r.base, r.cur)
	for _, cmp := range r.match {
		if cmp.f(record) != 0 {
			return false
		}
	}
	return true
}
