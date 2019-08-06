package path

import (
	"fmt"
)

type Path interface {
	String() string
	path()
}

type Value struct {
	Value string
}

type Term struct {
	Name string
}

type Inverse struct {
	Path Path
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

func (Value) path()        {}
func (Term) path()         {}
func (Inverse) path()      {}
func (Join) path()         {}
func (Intersection) path() {}
func (Union) path()        {}

func (p Value) String() string        { return fmt.Sprintf("'%s'", p.Value) }
func (p Term) String() string         { return p.Name }
func (p Inverse) String() string      { return fmt.Sprintf("~(%s)", p.Path) }
func (p Join) String() string         { return fmt.Sprintf("(%s)/(%s)", p.Left, p.Right) }
func (p Intersection) String() string { return fmt.Sprintf("(%s)&(%s)", p.Left, p.Right) }
func (p Union) String() string        { return fmt.Sprintf("(%s)|(%s)", p.Left, p.Right) }

type Set interface {
	Inverse() Set
	Join(Set) Set
	Union(Set) Set
	Intersection(Set) Set
}

type Env interface {
	Lookup(name string) (Set, error)
	Wrap(value string) (Set, error)
}

func Eval(path Path, env Env) (Set, error) {
	switch path := path.(type) {
	case Value:
		return env.Wrap(path.Value)
	case Term:
		return env.Lookup(path.Name)
	case Inverse:
		s, err := Eval(path.Path, env)
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
