package path

import (
	"testing"
)

func TestAnalysis(t *testing.T) {
	p, _ := Parse([]byte(`b_name/~name&~a/b`))
	pt, _ := Eval(p, testEnv{})
	if pt != testSet(`b_name/~name&~a/b`) {
		t.Fail()
	}
}

type testEnv map[string]Set

func (e testEnv) Lookup(name string) (Set, error) {
	return testSet(name), nil
}

func (e testEnv) Wrap(value string) (Set, error) {
	return testSet("'" + value + "'"), nil
}

type testSet string

func (s testSet) Inverse() Set {
	return "~" + s
}

func (s testSet) Union(t Set) Set {
	return s + "|" + t.(testSet)
}

func (s testSet) Intersection(t Set) Set {
	return s + "&" + t.(testSet)
}

func (s testSet) Join(t Set) Set {
	return s + "/" + t.(testSet)
}
