package er

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

func TestQuery(t *testing.T) {
	type Record struct {
		unexported string
		ParentID   int    `key:"PK:0"`
		Name       string `key:"PK:1"`
		Age        int
	}
	columns := SpecFor(reflect.TypeOf(Record{}))
	table := []Record{
		{
			"private",
			0,
			"A",
			1,
		},
		{
			"private",
			0,
			"B",
			2,
		},
		{
			"private",
			0,
			"C",
			1,
		},
		{
			"private",
			1,
			"A",
			3,
		},
	}
	runQuery := func(q Query) []Record {
		r := columns.EvalQuery(q, table)
		var res []Record
		for r.Next() {
			res = append(res, table[r.This()])
		}
		return res
	}
	assert.Equal(t, table[:3], runQuery(Query{0}))
	assert.Equal(t, table[1:2], runQuery(Query{0, "B"}))
	assert.Equal(t, table[1:2], runQuery(Query{1: "B"}))
	assert.Equal(t, []Record{table[0], table[2]}, runQuery(Query{2: 1}))
}

type benchmarkRecord1k struct {
	ID int `key:"PK"`
}

type benchmarkRecord2k struct {
	ParentID int `key:"PK"`
	ID       int `key:"PK"`
}

var (
	benchmarkRes        *QueryResult
	benchmarkMeta1k     = SpecFor(reflect.TypeOf(benchmarkRecord1k{}))
	benchmarkTestData1k = func() []benchmarkRecord1k {
		res := make([]benchmarkRecord1k, 10000)
		for i := range res {
			res[i].ID = i
		}
		return res
	}()
	benchmarkMeta2k     = SpecFor(reflect.TypeOf(benchmarkRecord2k{}))
	benchmarkTestData2k = func() []benchmarkRecord2k {
		res := make([]benchmarkRecord2k, 10000)
		for i := range res {
			res[i].ParentID = i / 100
			res[i].ID = i
		}
		return res
	}()
)

func BenchmarkIndex1k1r(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{0}, benchmarkTestData1k[:1])
	}
}

func BenchmarkIndex1k10r(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{0}, benchmarkTestData1k[:10])
	}
}

func BenchmarkIndex1k100r(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{0}, benchmarkTestData1k[:100])
	}
}

func BenchmarkIndex1k1000r(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{0}, benchmarkTestData1k[:1000])
	}
}

func BenchmarkIndex1k10000r(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{0}, benchmarkTestData1k[:10000])
	}
}

func BenchmarkRandomIndex1k1r(b *testing.B) {
	n := 1
	for i := 0; i < b.N; i++ {
		k := rand.Intn(n)
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{k}, benchmarkTestData1k[:n])
	}
}

func BenchmarkRandomIndex1k10r(b *testing.B) {
	n := 10
	for i := 0; i < b.N; i++ {
		k := rand.Intn(n)
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{k}, benchmarkTestData1k[:n])
	}
}

func BenchmarkRandomIndex1k100r(b *testing.B) {
	n := 100
	for i := 0; i < b.N; i++ {
		k := rand.Intn(n)
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{k}, benchmarkTestData1k[:n])
	}
}

func BenchmarkRandomIndex1k1000r(b *testing.B) {
	n := 1000
	for i := 0; i < b.N; i++ {
		k := rand.Intn(n)
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{k}, benchmarkTestData1k[:n])
	}
}

func BenchmarkRandomIndex1k10000r(b *testing.B) {
	n := 10000
	for i := 0; i < b.N; i++ {
		k := rand.Intn(n)
		benchmarkRes = benchmarkMeta1k.EvalQuery(Query{k}, benchmarkTestData1k[:n])
	}
}
