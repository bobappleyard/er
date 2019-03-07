package er

import (
	"github.com/bobappleyard/er"
	"github.com/bobappleyard/top"
	"github.com/bobappleyard/unify"
)

func sortRels(m *er.EntityModel) ([]*er.Relationship, error) {
	var g top.Graph
	for _, t := range m.Types {
		g.Link(m, t)
		for _, r := range t.Relationships {
			g.Link(r.Target, r)
			for _, c := range r.Constraints {
				for _, c := range c.Diagonal.Components {
					g.Link(c.Rel, r)
				}
				for _, c := range c.Riser.Components {
					g.Link(c.Rel, r)
				}
			}
			if r.Identifying {
				g.Link(r, t)
			}
		}
	}
	ires, err := g.Sort()
	if err != nil {
		return nil, err
	}
	res := make([]*er.Relationship, 0, len(ires))
	for _, r := range ires {
		if r, ok := r.(*er.Relationship); ok {
			res = append(res, r)
		}
	}
	return res, nil
}

func logicalToPhysical(m *er.EntityModel) error {
	rs, err := sortRels(m)
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Implementation != nil {
			continue
		}
		if err := (&relationshipImplementation{r: r}).create(); err != nil {
			return err
		}
	}
	return nil
}

type relationshipImplementation struct {
	r    *er.Relationship
	key  []unify.Var
	subs unify.Subs
}

func (r *relationshipImplementation) create() error {
	err := r.initProblem()
	if err != nil {
		return err
	}
	for _, c := range r.r.Constraints {
		err = r.applyConstraint(c)
		if err != nil {
			return err
		}
	}
	r.implement()
	return nil
}

func term(e *er.EntityType, f func(*er.Attribute) unify.Term) unify.Term {
	t := unify.Apply{Fn: e}
	for _, a := range e.Attributes {
		if !a.Identifying {
			continue
		}
		t.Args = append(t.Args, f(a))
	}
	return t
}

func (r *relationshipImplementation) initProblem() error {
	var key []unify.Var
	target := term(r.r.Target, func(a *er.Attribute) unify.Term {
		v := unify.Var{Of: a}
		key = append(key, v)
		return v
	})
	source := term(r.r.Target, func(a *er.Attribute) unify.Term {
		return unify.Var{Of: &er.Attribute{
			Owner:       r.r.Source,
			Name:        r.r.Name + "_" + a.Name,
			Type:        a.Type,
			Identifying: r.r.Identifying,
		}}
	})
	subs, err := unify.Unify(target, source, nil)
	r.key = key
	r.subs = subs
	return err
}

func (r *relationshipImplementation) applyConstraint(c er.Constraint) error {
	var source, target unify.Term
	var err error
	target, r.subs, err = followRiser(c.Riser, r.subs)
	if err != nil {
		return err
	}
	source, r.subs, err = followDiagonal(c.Diagonal, r.subs)
	if err != nil {
		return err
	}
	r.subs, err = unify.Unify(target, source, r.subs)
	return err
}

func (r *relationshipImplementation) implement() {
	r.r.Implementation = make([]er.Implementation, len(r.key))
	for i, a := range r.key {
		s := r.subs[a]
		targ := a.Of.(*er.Attribute)
		if s, ok := s.(unify.Var); ok {
			attr := s.Of.(*er.Attribute)
			r.r.Source.Attributes = append(r.r.Source.Attributes, attr)
			r.r.Implementation[i] = er.Implementation{
				Target: targ,
				Source: attr,
			}
			continue
		}
		src := s.(unify.Apply)
		r.r.Implementation[i] = er.Implementation{
			Target:   targ,
			Source:   src.Fn.(*er.Attribute),
			BasePath: pathFromTrace((src.Args[0].(unify.Apply)).Fn),
		}
	}
}

func followDiagonal(d er.Diagonal, subs unify.Subs) (unify.Term, unify.Subs, error) {
	f := func(a *er.Attribute, path []er.Component) unify.Term {
		return unify.Apply{Fn: a, Args: []unify.Term{unify.Apply{Fn: traceFromPath(path)}}}
	}
	return followPath(d.Components, subs, f)
}

func followRiser(r er.Riser, subs unify.Subs) (unify.Term, unify.Subs, error) {
	f := func(a *er.Attribute, path []er.Component) unify.Term {
		return unify.Var{Of: a}
	}
	return followPath(r.Components, subs, f)
}

func followPath(path []er.Component, subs unify.Subs, sourceAttr func(*er.Attribute, []er.Component) unify.Term) (unify.Term, unify.Subs, error) {
	var err error
	var target unify.Term
	for idx, cc := range path {
		cr := cc.Rel
		source := unify.Apply{Fn: cr.Target, Args: make([]unify.Term, len(cr.Implementation))}
		for i, a := range cr.Implementation {
			source.Args[i] = sourceAttr(a.Source, path[:idx])
		}
		target = term(cr.Target, func(a *er.Attribute) unify.Term {
			return unify.Var{Of: a}
		})
		subs, err = unify.Unify(source, target, subs)
		if err != nil {
			return nil, nil, err
		}
	}
	return target, subs, nil
}

type trace struct {
	cur  *er.Relationship
	next interface{}
}

func traceFromPath(p []er.Component) interface{} {
	if len(p) == 0 {
		return nil
	}
	return trace{
		p[0].Rel,
		traceFromPath(p[1:]),
	}
}

func pathFromTrace(t interface{}) []er.Component {
	var res []er.Component
	for {
		if t == nil {
			break
		}
		c := t.(trace)
		res = append(res, er.Component{Rel: c.cur})
		t = c.next
	}
	return res
}
