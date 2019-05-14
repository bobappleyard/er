package diagram

import (
	"github.com/bobappleyard/er"
	"math"
)

// tower represents a set of entities connected by dependency.
type tower struct {
	t                 *er.EntityType
	ts                map[*er.EntityType]*tower
	head, body        box
	up                *tower
	left, right, down towers
	lines             []segment
}

// towers is a list of towers
type towers []*tower

// box represents the spatial extent of part of a tower.
type box struct {
	x, y, w, h int
}

// link represents an association between two towers
type link struct {
	from, to *tower
}

type segment struct {
	x1, y1, x2, y2 int
}

func buildTowers(m *er.EntityModel) *tower {
	ts := map[*er.EntityType]*tower{}
	for _, t := range m.Types {
		ts[t] = &tower{t: t, ts: ts}
	}
	root := &tower{}
	for _, t := range m.Types {
		child := ts[t]
		parent := root
		if t.DependsOn != nil {
			parent = ts[t.DependsOn.Target]
		}
		parent.down = append(parent.down, child)
		child.up = parent
	}
	return root
}

func (t *tower) layoutDiagram() {
	t.calcEntityOrder()
	t.layoutBoxes(0, 0)
	t.layoutLines()
}

func (t *tower) calcEntityOrder() {
	links := t.linkMap()
	score := math.MaxInt32
	order := firstPerm(len(t.down))
	for p := firstPerm(len(t.down)); p != nil; p = nextPerm(p) {
		newScore := t.permutationCost(p, links)
		if newScore < score {
			score = newScore
			copy(order, p)
		}
	}
	newOrder := make(towers, len(t.down))
	for i, u := range t.down {
		newOrder[order[i]] = u
	}
	t.down = newOrder
	for i, u := range t.down {
		u.left = append(append([]*tower{}, t.left...), t.down[:i]...)
		u.right = append(append([]*tower{}, t.right...), t.down[i+1:]...)
		u.calcEntityOrder()
	}
}

func (t *tower) layoutBoxes(x, y int) {
	t.head.x = x
	t.head.y = y
	if t.t != nil {
		t.calcHeadBounds()
		y += t.head.h + 10
	}
	for _, t := range t.down {
		t.layoutBoxes(x, y)
		x += t.body.w + 10
	}
	t.calcBodyBounds()
	t.adjustHeadBounds()
}

func (t *tower) layoutLines() {
	if t.t != nil {
		x1, y1 := t.head.center()
		for _, r := range t.t.Relationships {
			targ := t.ts[r.Target]
			if r == t.t.DependsOn {
				x1, y1 := x1, t.head.y
				x2, y2 := x1, targ.head.y+targ.head.h
				t.lines = append(t.lines,
					segment{x1, y1 + 5, x2, y2},
					segment{x1 - 3, y1 + 5, x1, y1},
					segment{x1 + 3, y1 + 5, x1, y1},
				)
				if r.Identifying {
					t.lines = append(t.lines, segment{x1 - 3, y1 - 2, x1 + 3, y1 - 2})
				}
				continue
			}
			if targ.head.x < t.head.x {
				_, cy := targ.head.center()
				x1, y1 := t.head.x, y1
				x2, y2 := targ.head.x+targ.head.w, cy
				midpoint := x2 + (x1-x2)/2
				t.lines = append(t.lines,
					segment{x1 + 5, y1, midpoint, y1},
					segment{midpoint, y1, midpoint, y2},
					segment{midpoint, y2, x2, y2},
					segment{x1 + 5, y1 - 3, x1, y1},
					segment{x1 + 5, y1 + 3, x1, y1},
				)
				if r.Identifying {
					t.lines = append(t.lines, segment{x1 - 2, y1 - 3, x1 - 2, y1 + 3})
				}
			}
		}
	}
	for _, t := range t.down {
		t.layoutLines()
	}
}

func (b box) center() (x, y int) {
	return b.x + (b.w / 2), b.y + (b.h / 2)
}

var (
	left  = &tower{}
	right = &tower{}
)

func (t *tower) linkMap() map[link]int {
	res := map[link]int{}
	addWeight := func(t, u *tower, v towers) {
		weight := t.connections(v)
		if weight == 0 {
			return
		}
		res[link{t, u}] = weight
		res[link{u, t}] = weight
	}
	for i, u := range t.down {
		ut := towers{u}
		for _, v := range t.down[i+1:] {
			addWeight(u, v, ut)
		}
		addWeight(u, left, t.left)
		addWeight(u, right, t.right)
	}
	return res
}

func (t *tower) permutationCost(p []int, links map[link]int) int {
	res := 0
	scoreComponent := func(w, d int) {
		res += w * d * d
	}
	for i := range t.down {
		u := t.down[p[i]]
		for d := range t.down[i+1:] {
			v := t.down[p[d+i]]
			scoreComponent(links[link{u, v}], d)
		}
		scoreComponent(links[link{left, u}], i)
		scoreComponent(links[link{right, u}], len(p)-i-1)
	}
	return res
}

func (ts towers) contains(e *er.EntityType) bool {
	for _, t := range ts {
		if t.t == e {
			return true
		}
		if t.down.contains(e) {
			return true
		}
	}
	return false
}

func (t *tower) connections(u towers) int {
	c := 0
	for _, r := range t.t.Relationships {
		if u.contains(r.Target) {
			c++
		}
	}
	for _, t := range t.down {
		c += t.connections(u)
	}
	return c
}

func (t *tower) calcHeadBounds() {
	var w, h int
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
	t.head.w = w + 10
	t.head.h = h + 10
}

func (t *tower) calcBodyBounds() {
	t.body = t.head
	w, h := 0, 0
	for _, u := range t.down {
		if u.body.h > h {
			h = u.body.h
		}
		w += u.body.w
	}
	t.body.h += h + 10
	if w > t.body.w {
		t.body.w = w
	}
	t.body.w += 10
}

func (t *tower) adjustHeadBounds() {
	switch len(t.down) {
	case 0:
	case 1:
		x, _ := t.down[0].head.center()
		t.head.x = x - t.head.w/2
	default:
		fx, _ := t.down[0].head.center()
		lx, _ := t.down[len(t.down)-1].head.center()
		if lx-fx > t.head.w && (lx-fx+20) < t.body.w {
			t.head.x = fx - 10
			t.head.w = lx - fx + 20
		}
	}
}
