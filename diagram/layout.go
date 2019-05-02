package diagram

import (
	"github.com/bobappleyard/er"
)

// tower represents a set of entities connected by dependency.
type tower struct {
	t  *er.EntityType
	cs []*tower
}

func towersFor(m *er.EntityModel) *tower {
	ts := map[*er.EntityType]*tower{}
	for _, t := range m.Types {
		ts[t] = &tower{t: t}
	}
	var roots []*tower
	for _, t := range m.Types {
		if t.DependsOn == nil {
			roots = append(roots, ts[t])
			continue
		}
		parent := ts[t.DependsOn.Target]
		parent.cs = append(parent.cs, ts[t])
	}
	return &tower{cs: roots}
}

func (t *tower) bounds() (w, h int) {
	if t.t == nil {
		return 0, 0
	}
	h = 15 * (len(t.t.Attributes) + 1)
	for _, a := range t.t.Attributes {
		est := 15 + len(a.Name)*7
		if est > w {
			w = est
		}
	}
	if w < 50 {
		w = 50
	}
	if h < 35 {
		h = 35
	}
	return w + 10, h + 10
}

func (t *tower) depth() int {
	d := 0
	for _, t := range t.cs {
		e := t.depth()
		if e > d {
			d = e
		}
	}
	_, h := t.bounds()
	return d + h
}

func (t *tower) width() int {
	w, _ := t.bounds()
	if len(t.cs) == 0 {
		return w
	}
	d := 0
	for _, t := range t.cs {
		d += t.width()
	}
	if w > d {
		d = w
	}
	return d
}
