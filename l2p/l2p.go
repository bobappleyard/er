package l2p

import (
	"errors"
	"fmt"
	"github.com/bobappleyard/er"
	"github.com/bobappleyard/er/util/path"
	"io"
	"sort"
)

var (
	errModelChanged = errors.New("model changed")
	errModelLoop    = errors.New("model has a loop")
	ErrNoPath       = errors.New("no path through model")
)

func LogicalToPhysical(m *er.EntityModel) error {
	for {
		err := performAnalysis(modelEnv{m: m})
		if u, ok := err.(*modelUpdate); ok {
			u.r.Source.Attributes = append(u.r.Source.Attributes, u.attrs...)
			u.r.Path = u.path.String()
			continue
		}
		return err
	}
}

func performAnalysis(n modelEnv) error {
	for _, e := range n.m.Types {
		for _, r := range e.Relationships {
			_, err := n.implement(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type modelEnv struct {
	m       *er.EntityModel
	visited []*er.Relationship
}
type pathType []partition
type partition struct {
	source, target set
}
type set struct {
	entityType *er.EntityType
	attrs      uint64
}
type modelUpdate struct {
	r     *er.Relationship
	attrs []*er.Attribute
	path  path.Path
}

func (*modelUpdate) Error() string { return "update to model" }

var (
	absolute  = &er.EntityType{Name: "*"}
	attribute = &er.EntityType{Name: "$"}
)

func (n modelEnv) extend(r *er.Relationship) modelEnv {
	f := modelEnv{
		m:       n.m,
		visited: append([]*er.Relationship{}, n.visited...),
	}
	f.visited = append(f.visited, r)
	return f
}

func (n modelEnv) implement(r *er.Relationship) (path.Path, error) {
	for _, r := range r.Target.Relationships {
		if !r.Identifying {
			continue
		}
		_, err := n.extend(r).implement(r)
		if err != nil {
			return nil, err
		}
	}
	var err error
	var s path.Set = pathType{}
	p := path.Chart(
		path.InverseTerm{r.Source.Name},
		path.Term{"*"},
		path.Term{r.Target.Name},
	)
	if r.Path != "" {
		q, err := path.Parse([]byte(r.Path))
		if err != nil {
			return nil, err
		}
		p = path.Intersection{p, q}
	}
	s, err = path.Eval(p, n)
	if err != nil {
		return nil, err
	}
	return n.addMissingAttributes(r, p, s.(pathType))
}

func (n modelEnv) addMissingAttributes(r *er.Relationship, p path.Path, s pathType) (path.Path, error) {
	if len(s) == 0 {
		return nil, ErrNoPath
	}
	var target set
	for _, p := range s {
		if target.entityType == nil || p.target.attrCount() > target.attrCount() {
			target = p.target
		}
	}
	var err error
	var attrs []*er.Attribute
	for i, a := range r.Target.Attributes {
		if !a.Identifying || target.contains(i) {
			continue
		}
		err = errModelChanged
		name := r.Name + "_" + a.Name
		attrs = append(attrs, &er.Attribute{
			Name:        name,
			Type:        a.Type,
			Identifying: r.Identifying,
			Owner:       r.Source,
		})
		part := path.Join{
			path.Term{name},
			path.InverseTerm{a.Name},
		}
		if p == nil {
			p = part
		} else {
			p = path.Intersection{p, part}
		}
	}
	if len(attrs) != 0 {
		return nil, &modelUpdate{
			r:     r,
			attrs: attrs,
			path:  p,
		}
	}
	return p, err
}

func (n modelEnv) Lookup(name string) (path.Set, error) {
	var res pathType
	var err error
	if name == "*" {
		return pathType{{
			source: set{entityType: absolute},
			target: set{entityType: absolute},
		}}, nil
	}
	for _, e := range n.m.Types {
		if e.Name == name {
			res = append(res, partition{
				source: set{entityType: absolute},
				target: set{entityType: e},
			})
		}
		for i, a := range e.Attributes {
			if a.Name != name {
				continue
			}
			res = append(res, partition{
				source: set{entityType: e, attrs: 1 << uint64(i)},
				target: set{entityType: attribute},
			})
		}
		for _, r := range e.Relationships {
			if r.Name != name || n.seen(r) {
				continue
			}
			f := n.extend(r)
			p, e := f.implement(r)
			if e, ok := e.(*modelUpdate); ok {
				return nil, e
			}
			if e != nil {
				err = e
				continue
			}
			s, err := path.Eval(p, f)
			if err != nil {
				return nil, err
			}
			res = append(res, (s.(pathType))...)
		}
	}
	res = res.prune()
	if len(res) != 0 {
		return res, nil
	}
	return res, err
}

func (e modelEnv) seen(r *er.Relationship) bool {
	for _, s := range e.visited {
		if s == r {
			return true
		}
	}
	return false
}

func (pt pathType) Inverse() path.Set {
	res := make(pathType, len(pt))
	for i, p := range pt {
		res[i].source = p.target
		res[i].target = p.source
	}
	return res
}

func (left pathType) Join(right path.Set) path.Set {
	var res pathType
	for _, left := range left {
		for _, right := range right.(pathType) {
			if left.target.singleEntity().subsetOf(right.source) {
				res = append(res, partition{
					source: left.source,
					target: right.target,
				})
			}
		}
	}
	return res.prune()
}

func (left pathType) Intersection(right path.Set) path.Set {
	var res pathType
	for _, left := range left {
		for _, right := range right.(pathType) {
			if left.source.entityType != right.source.entityType {
				continue
			}
			if left.target.entityType != right.target.entityType {
				continue
			}
			res = append(res, left, right, partition{
				source: left.source.mergeWith(right.source).singleEntity(),
				target: left.target.mergeWith(right.target).singleEntity(),
			})
		}
	}
	return res.prune()
}

func (left pathType) Union(right path.Set) path.Set {
	panic("unimplemented")
}

func (pt pathType) prune() pathType {
	if len(pt) < 2 {
		return pt
	}
	sort.Slice(pt, func(i, j int) bool { return pt[i].lessThan(pt[j]) })
	res := make(pathType, 0, len(pt))
	res = append(res, pt[0])
	for i, q := range pt[1:] {
		p := pt[i]
		if p.lessThan(q) || q.lessThan(p) {
			res = append(res, q)
		}
	}
	return res
}

func (p partition) lessThan(q partition) bool {
	if p.source.entityType.Name != q.source.entityType.Name {
		return p.source.entityType.Name < q.source.entityType.Name
	}
	if p.target.entityType.Name != q.target.entityType.Name {
		return p.target.entityType.Name < q.target.entityType.Name
	}
	if p.source.attrs != q.source.attrs {
		return p.source.attrs < q.source.attrs
	}
	if p.target.attrs != q.target.attrs {
		return p.target.attrs < q.target.attrs
	}
	return false
}

func (left set) subsetOf(right set) bool {
	if left.entityType != right.entityType {
		return false
	}
	return left.attrs&right.attrs == right.attrs
}

func (s set) attrCount() int {
	v := s.attrs
	c := 0
	for v != 0 {
		c++
		v &= v - 1
	}
	return c
}

func (s set) contains(i int) bool {
	return s.attrs&(1<<uint64(i)) != 0
}

func (left set) mergeWith(right set) set {
	return set{left.entityType, left.attrs | right.attrs}
}

func (s set) singleEntity() set {
	for i, a := range s.entityType.Attributes {
		if !a.Identifying {
			continue
		}
		if !s.contains(i) {
			return s
		}
	}
	allAttrs := (1 << uint64(len(s.entityType.Attributes))) - 1
	return set{s.entityType, uint64(allAttrs)}
}

func (pt pathType) Format(t fmt.State, c rune) {
	for i, p := range pt {
		if i != 0 {
			io.WriteString(t, " | ")
		}
		fmt.Fprint(t, p)
	}
}

func (p partition) Format(t fmt.State, c rune) {
	fmt.Fprintf(t, "%s -> %s", p.source, p.target)
}

func (s set) Format(t fmt.State, c rune) {
	io.WriteString(t, s.entityType.Name)
	io.WriteString(t, "(")
	written := false
	for i, a := range s.entityType.Attributes {
		if s.contains(i) {
			if written {
				io.WriteString(t, ",")
			}
			written = true
			io.WriteString(t, a.Name)
		}
	}
	io.WriteString(t, ")")
}
