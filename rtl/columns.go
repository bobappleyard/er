package rtl

type StringColumn struct {
	ColumnID int
	Key      bool
	Val      func(idx int) string
}

type IntColumn struct {
	ColumnID int
	Key      bool
	Val      func(idx int) int
}
