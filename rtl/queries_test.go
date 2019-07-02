package rtl

import (
	"reflect"
	"testing"
)

func TestQueryConstruction(t *testing.T) {
	calls := -1
	testQuery := func() Query {
		calls++
		return queryForClause(clause{
			columnID: calls,
		})
	}
	assertClauses := func(q Query, ids ...int) func(*testing.T) {
		return func(t *testing.T) {
			if len(q.clauses) != len(ids) {
				t.Errorf("got %d clauses, expected %d", len(q.clauses), len(ids))
				return
			}
			for i, id := range ids {
				if q.clauses[i].columnID != id {
					t.Errorf("[%d] got clause %d, expecting %d", i, q.clauses[i].columnID, id)
				}
			}
		}
	}
	q := testQuery()
	r := testQuery()
	s := q.And(r)
	t.Run("q1", assertClauses(q, 0))
	t.Run("r1", assertClauses(r, 1))
	t.Run("s1", assertClauses(s, 0, 1))
	u := testQuery()
	v := q.And(u)
	t.Run("q2", assertClauses(q, 0))
	t.Run("u2", assertClauses(u, 2))
	t.Run("v2", assertClauses(v, 0, 2))
	t.Run("s2", assertClauses(s, 0, 1))
	w := q.Or(u)
	t.Run("w1", assertClauses(w, 0))
	if w.alt == nil {
		t.Error("missing alt")
	}
	t.Run("w2", assertClauses(*w.alt, 2))
}

func TestQueryEvaluation(t *testing.T) {
	type dept struct {
		orderID   int
		productID string
		quantity  int
	}
	rows := []dept{
		{1, "Apple", 10},
		{1, "Banana", 20},
		{1, "Butter", 1},
		{2, "Apple", 10},
		{2, "Cider", 1},
	}
	orderID := IntIndex(0, func(idx int) int { return rows[idx].orderID })
	productID := StringIndex(1, func(idx int) string { return rows[idx].productID })
	quantity := IntColumn(2, func(idx int) int { return rows[idx].quantity })

	runQuery := func(q Query) (res []dept) {
		r := EvalQuery(q, len(rows))
		for r.Next() {
			res = append(res, rows[r.This()])
		}
		return res
	}
	for _, test := range []struct {
		name string
		q    Query
		rows []dept
	}{
		{
			name: "ByID",
			q:    orderID.Eq(1),
			rows: []dept{
				{1, "Apple", 10},
				{1, "Banana", 20},
				{1, "Butter", 1},
			},
		},
		{
			name: "ByIdAndName",
			q:    orderID.Eq(1).And(productID.Eq("Apple")),
			rows: []dept{
				{1, "Apple", 10},
			},
		},
		{
			name: "LargeQuantity",
			q:    quantity.Gt(1),
			rows: []dept{
				{1, "Apple", 10},
				{1, "Banana", 20},
				{2, "Apple", 10},
			},
		},
		{
			name: "NameOrQuantity",
			q:    productID.Eq("Apple").Or(quantity.Eq(20)),
			rows: []dept{
				{1, "Apple", 10},
				{1, "Banana", 20},
				{2, "Apple", 10},
			},
		},
		{
			name: "AppleOrder1OrQuantity",
			q:    productID.Eq("Apple").And(orderID.Eq(1)).Or(quantity.Eq(20)),
			rows: []dept{
				{1, "Apple", 10},
				{1, "Banana", 20},
			},
		},
		{
			name: "ThreeAlts",
			q:    productID.Eq("Apple").Or(productID.Eq("Banana")).Or(productID.Eq("Butter")),
			rows: []dept{
				{1, "Apple", 10},
				{1, "Banana", 20},
				{1, "Butter", 1},
				{2, "Apple", 10},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := runQuery(test.q)
			if !reflect.DeepEqual(got, test.rows) {
				t.Errorf("got %v, expected %v", got, test.rows)
			}
		})
	}
}
