package er

import (
	"github.com/bobappleyard/top"
	"github.com/bobappleyard/unify"
)

func (m *EntityModel) sortRels() ([]*Relationship, error) {
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
	res := make([]*Relationship, 0, len(ires))
	for _, r := range ires {
		if r, ok := r.(*Relationship); ok {
			res = append(res, r)
		}
	}
	return res, nil
}

func (m *EntityModel) logicalToPhysical() error {
	rs, err := m.sortRels()
	if err != nil {
		return err
	}
	for _, r := range rs {
		err := r.logicalToPhysical()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EntityType) term(f func(*Attribute) unify.Term) unify.Term {
	t := unify.Apply{Fn: e}
	for _, a := range e.Attributes {
		if !a.Identifying {
			continue
		}
		t.Args = append(t.Args, f(a))
	}
	return t
}

func (r *Relationship) logicalToPhysical() error {
	key, subs, err := r.initProblem()
	if err != nil {
		return err
	}
	for _, c := range r.Constraints {
		subs, err = r.applyConstraint(c, subs)
		if err != nil {
			return err
		}
	}
	r.implement(key, subs)
	return nil
}

func (r *Relationship) initProblem() ([]unify.Var, unify.Subs, error) {
	var key []unify.Var
	target := r.Target.term(func(a *Attribute) unify.Term {
		v := unify.Var{Of: a}
		key = append(key, v)
		return v
	})
	source := r.Target.term(func(a *Attribute) unify.Term {
		return unify.Var{Of: &Attribute{
			Owner:       r.Source,
			Name:        r.Name + "_" + a.Name,
			Type:        a.Type,
			Identifying: r.Identifying,
		}}
	})
	subs, err := unify.Unify(target, source, nil)
	return key, subs, err
}

func (r *Relationship) applyConstraint(c Constraint, subs unify.Subs) (unify.Subs, error) {
	var source, target unify.Term
	var err error
	target, subs, err = followPath{
		path: c.Riser.Components,
		subs: subs,
		sourceAttr: func(a *Attribute, path []Component) unify.Term {
			return unify.Var{Of: a}
		},
	}.eval()
	if err != nil {
		return nil, err
	}
	source, subs, err = followPath{
		path: c.Diagonal.Components,
		subs: subs,
		sourceAttr: func(a *Attribute, path []Component) unify.Term {
			return unify.Apply{Fn: a, Args: []unify.Term{unify.Apply{Fn: traceFromPath(path)}}}
		},
	}.eval()
	if err != nil {
		return nil, err
	}
	subs, err = unify.Unify(target, source, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *Relationship) implement(key []unify.Var, subs unify.Subs) {
	r.Implementation = make([]Implementation, len(key))
	for i, a := range key {
		s := subs[a]
		targ := a.Of.(*Attribute)
		if s, ok := s.(unify.Var); ok {
			attr := s.Of.(*Attribute)
			r.Source.Attributes = append(r.Source.Attributes, attr)
			r.Implementation[i] = Implementation{
				Target: targ,
				Source: attr,
			}
			continue
		}
		src := s.(unify.Apply)
		r.Implementation[i] = Implementation{
			Target:   targ,
			Source:   src.Fn.(*Attribute),
			BasePath: pathFromTrace((src.Args[0].(unify.Apply)).Fn),
		}
	}
}

type followPath struct {
	path       []Component
	subs       unify.Subs
	sourceAttr func(*Attribute, []Component) unify.Term
}

func (p followPath) eval() (unify.Term, unify.Subs, error) {
	var err error
	var target unify.Term
	subs := p.subs
	for idx, cc := range p.path {
		cr := cc.Rel
		source := unify.Apply{Fn: cr.Target, Args: make([]unify.Term, len(cr.Implementation))}
		for i, a := range cr.Implementation {
			source.Args[i] = p.sourceAttr(a.Source, p.path[:idx])
		}
		target = cr.Target.term(func(a *Attribute) unify.Term {
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
	cur  *Relationship
	next interface{}
}

func traceFromPath(p []Component) interface{} {
	if len(p) == 0 {
		return nil
	}
	return trace{
		p[0].Rel,
		traceFromPath(p[1:]),
	}
}

func pathFromTrace(t interface{}) []Component {
	var res []Component
	for {
		if t == nil {
			break
		}
		c := t.(trace)
		res = append(res, Component{Rel: c.cur})
		t = c.next
	}
	return res
}
