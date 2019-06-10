package rtl

import (
	"sort"
)

type Query struct {
	clauses []clause
	alt     *Query
}

type QueryResult struct {
	this, end int
	q         Query
}

type test byte

const (
	eq test = iota
	lt
	le
	gt
	ge
	ne
	key
)

type clause struct {
	columnID int
	op       test
	cmp      func(int) int
}

// Query construction

func queryForClause(c clause) Query {
	return Query{clauses: []clause{c}}
}

// Query composition

func (q Query) And(r Query) Query {
	var res Query
	res.clauses = append(res.clauses, q.clauses...)
	res.clauses = append(res.clauses, r.clauses...)
	if q.alt != nil {
		alt := q.alt.And(r)
		res.alt = &alt
	}
	return res
}

func (q Query) Or(r Query) Query {
	return Query{
		clauses: q.clauses,
		alt:     &r,
	}
}

// Query evaluation

func All(n int) *QueryResult {
	return EvalQuery(queryForClause(clause{
		cmp: func(int) int { return 0 },
	}), n)
}

func EvalQuery(q Query, n int) *QueryResult {
	res := &QueryResult{0, n, q}
	res.refineSearchSpace()
	res.this--
	return res
}

func (c clause) key() bool { return c.op == key }

func (r *QueryResult) refineSearchSpace() {
	this := r.this
	end := r.end
	for q := &r.q; q != nil; q = q.alt {
		min := r.this
		max := r.end
		clauses := q.clauses
		sort.Slice(clauses, func(i, j int) bool {
			if clauses[i].key() != clauses[j].key() {
				return clauses[i].key()
			}
			return clauses[i].columnID < clauses[j].columnID
		})
		for i, c := range clauses {
			if c.op != key || i != c.columnID {
				break
			}
			min = sort.Search(max-min, func(idx int) bool {
				return c.cmp(idx+min) >= 0
			}) + min
			max = sort.Search(max-min, func(idx int) bool {
				return c.cmp(idx+min) > 0
			}) + min
		}
		if min < this || q == &r.q {
			this = min
		}
		if max > end || q == &r.q {
			end = max
		}
	}
	r.this = this
	r.end = end
}

func (r *QueryResult) This() int {
	return r.this
}

func (r *QueryResult) Next() bool {
	for {
		r.this++
		if r.this >= r.end {
			return false
		}
		for q := &r.q; q != nil; q = q.alt {
			if q.matches(r.this) {
				return true
			}
		}
	}
}

func (q Query) matches(idx int) bool {
	for _, c := range q.clauses {
		if !c.matches(idx) {
			return false
		}
	}
	return true
}

func (c clause) matches(idx int) bool {
	cmp := c.cmp(idx)
	switch c.op {
	case key, eq:
		return cmp == 0
	case lt:
		return cmp < 0
	case le:
		return cmp <= 0
	case gt:
		return cmp > 0
	case ge:
		return cmp >= 0
	case ne:
		return cmp != 0
	}
	panic("invalid op")
}
