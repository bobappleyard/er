package path

import (
	"fmt"
)

type Path interface {
	String() string
	path()
}

type Term struct {
	Name string
}

type InverseTerm struct {
	Name string
}

type Join struct {
	Left, Right Path
}

type Intersection struct {
	Left, Right Path
}

type Union struct {
	Left, Right Path
}

func (Term) path()         {}
func (InverseTerm) path()  {}
func (Join) path()         {}
func (Intersection) path() {}
func (Union) path()        {}

func (p Term) String() string         { return p.Name }
func (p InverseTerm) String() string  { return "~" + p.Name }
func (p Join) String() string         { return fmt.Sprintf("(%s)/(%s)", p.Left, p.Right) }
func (p Intersection) String() string { return fmt.Sprintf("(%s)&(%s)", p.Left, p.Right) }
func (p Union) String() string        { return fmt.Sprintf("(%s)|(%s)", p.Left, p.Right) }

func Chart(ps ...Path) Path {
	res := ps[len(ps)-1]
	for i := len(ps) - 2; i >= 0; i-- {
		res = Join{ps[i], res}
	}
	return res
}

func Inverse(path Path) Path {
	switch path := path.(type) {
	case Term:
		return InverseTerm{path.Name}
	case InverseTerm:
		return Term{path.Name}
	case Join:
		return Join{Inverse(path.Right), Inverse(path.Left)}
	case Intersection:
		return Intersection{Inverse(path.Left), Inverse(path.Right)}
	case Union:
		return Union{Inverse(path.Left), Inverse(path.Right)}
	}
	panic("unreachable")
}

type Set interface {
	Inverse() Set
	Join(Set) Set
	Union(Set) Set
	Intersection(Set) Set
}

type Env interface {
	Lookup(name string) (Set, error)
}

func Eval(path Path, env Env) (Set, error) {
	switch path := path.(type) {
	case Term:
		return env.Lookup(path.Name)
	case InverseTerm:
		s, err := env.Lookup(path.Name)
		if err != nil {
			return nil, err
		}
		return s.Inverse(), nil
	case Join:
		return evalPair(path.Left, path.Right, env, Set.Join)
	case Intersection:
		return evalPair(path.Left, path.Right, env, Set.Intersection)
	case Union:
		return evalPair(path.Left, path.Right, env, Set.Union)
	}
	panic("unreachable")
}

func evalPair(left, right Path, env Env, f func(Set, Set) Set) (Set, error) {
	l, err := Eval(left, env)
	if err != nil {
		return nil, err
	}
	r, err := Eval(right, env)
	if err != nil {
		return nil, err
	}
	return f(l, r), nil
}
