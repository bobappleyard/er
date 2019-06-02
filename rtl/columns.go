package rtl

import (
	"sort"
)

type StringColumn struct {
	data  []string
	count int
}

type StringIndex struct {
	StringColumn
}

func (c *StringColumn) Insert(value string) {
	c.data = append(c.data, value)
}

func (c *StringColumn) len() int {
	return len(c.data)
}

func (c *StringColumn) cmp(i, j int) int {
	return 0
}

func (c *StringColumn) apply(idxs []int) {
	newData := make([]string, len(c.data))
	for i, s := range c.data {
		newData[idxs[i]] = s
	}
	c.data = newData
}

func (c *StringIndex) cmp(i, j int) int {
	a, b := c.data[i], c.data[j]
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

type column interface {
	len() int
	cmp(i, j int) int
	apply(idxs []int)
}

func EnsureUniqueness(cs ...column) error {
	n := cs[0].len()
	idxs := make([]int, n)
	for i := range idxs {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i, j int) bool {
		i = idxs[i]
		j = idxs[j]
		for _, c := range cs {
			cmp := c.cmp(i, j)
			if cmp < 0 {
				return true
			}
			if cmp > 0 {
				return false
			}
		}
		return true
	})
	for _, c := range cs {
		c.apply(idxs)
	}
	return nil
}
