package l2p

import (
	"github.com/bobappleyard/er"
	"github.com/bobappleyard/er/util/path"
)

// Apply a very basic type checking/inference algorithm to the paths that
// implement the provided relationship.
func resolveRelationship(m *er.EntityModel, r *er.Relationship) (resolvedPath, error) {
	var p path.Path = path.Join{
		Left: path.Inverse{Path: path.Term{Name: r.Source.Name}},
		Right: path.Join{
			Left:  path.Term{Name: "*"},
			Right: path.Term{Name: r.Target.Name},
		},
	}
	if r.Path != "" {
		q, err := path.Parse([]byte(r.Path))
		if err != nil {
			return nil, err
		}
		p = path.Intersection{Left: p, Right: q}
	}
	s, err := path.Eval(p, resolveEnv{m})
	if err != nil {
		return nil, err
	}
	rp := s.(potentialPathType)
	if len(rp) == 0 {
		return nil, ErrNoPath
	}
	if len(rp) > 1 {
		return nil, ErrAmbiguousPath
	}
	return rp[0], nil
}

type resolveEnv struct {
	m *er.EntityModel
}

type potentialPathType []resolvedPath

func (e resolveEnv) Lookup(name string) (path.Set, error) {
	if name == "*" {
		return potentialPathType{absolute{}}, nil
	}
	var res potentialPathType
	for _, e := range e.m.Types {
		if e.Name == name {
			res = append(res, resolvedEntityType{e})
		}
		for _, a := range e.Attributes {
			if a.Name == name {
				res = append(res, resolvedAttribute{a})
			}
		}
		for _, r := range e.Relationships {
			if r.Name == name {
				res = append(res, resolvedRelationship{r})
			}
		}
	}
	return res, nil
}

func (e resolveEnv) Wrap(value string) (path.Set, error) {
	return potentialPathType{resolvedValue{value}}, nil
}

func (t potentialPathType) Inverse() path.Set {
	res := make(potentialPathType, len(t))
	for i, r := range t {
		res[i] = resolvedInverse{r}
	}
	return res
}

func (t potentialPathType) Join(u path.Set) path.Set {
	var res potentialPathType
	for _, t := range t {
		_, target := t.route()
		for _, u := range u.(potentialPathType) {
			source, _ := u.route()
			if source != target {
				continue
			}
			res = append(res, resolvedJoin{t, u})
		}
	}
	return res
}

func (t potentialPathType) Intersection(u path.Set) path.Set {
	var res potentialPathType
	for _, t := range t {
		sourceT, targetT := t.route()
		for _, u := range u.(potentialPathType) {
			sourceU, targetU := u.route()
			if sourceT != sourceU {
				continue
			}
			if targetT != targetU {
				continue
			}
			res = append(res, resolvedIntersection{t, u})
		}
	}
	return res
}

func (t potentialPathType) Union(u path.Set) path.Set {
	panic("unsupported")
}

func (absolute) route() (source, target *er.EntityType) {
	return absoluteType, absoluteType
}

func (p resolvedValue) route() (source, target *er.EntityType) {
	return absoluteType, valueType
}

func (p resolvedEntityType) route() (source, target *er.EntityType) {
	return absoluteType, p.e
}

func (p resolvedRelationship) route() (source, target *er.EntityType) {
	return p.r.Source, p.r.Target
}

func (p resolvedAttribute) route() (source, target *er.EntityType) {
	return p.a.Owner, valueType
}

func (p resolvedInverse) route() (source, target *er.EntityType) {
	target, source = p.p.route()
	return source, target
}

func (p resolvedJoin) route() (source, target *er.EntityType) {
	source, _ = p.left.route()
	_, target = p.right.route()
	return source, target
}

func (p resolvedIntersection) route() (source, target *er.EntityType) {
	return p.left.route()
}
