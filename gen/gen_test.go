package gen

import (
	. "github.com/bobappleyard/er"
	"github.com/bobappleyard/er/l2p"
	"io/ioutil"
	"os/exec"
	"path"
	"testing"
)

func TestGen(t *testing.T) {
	m := &EntityModel{
		Name: "square",
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
		},
	}
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	d := m.Types[3]
	for _, t := range m.Types {
		t.Attributes = []*Attribute{
			{
				Owner:       t,
				Name:        "name",
				Type:        StringType,
				Identifying: true,
			},
		}
	}
	a.Relationships = []*Relationship{
		{
			Name:   "s",
			Source: a,
			Target: b,
		},
	}
	c.Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: c,
			Target: a,
		},
		{
			Name:   "f",
			Source: c,
			Target: d,
		},
	}
	d.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      d,
			Target:      b,
			Identifying: true,
		},
	}
	c.Relationships[1].Constraints = []Constraint{
		{
			Diagonal{[]Component{
				{Rel: c.Relationships[0]},
				{Rel: a.Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: d.Relationships[0]},
			}},
		},
	}
	c.DependsOn = c.Relationships[0]
	d.DependsOn = d.Relationships[0]

	l2p.LogicalToPhysical(m)
	bs, err := generate(m)
	if err != nil {
		t.Error(err)
		return
	}
	ioutil.WriteFile(path.Join("test", "pkg.go"), bs, 0777)
	cmd := exec.Command("go", "build", "-o", "test/pkg.o", "test/pkg.go")
	bs, err = cmd.CombinedOutput()
	cmd = exec.Command("go", "test", "./test")
	bs, err = cmd.CombinedOutput()
	t.Log(string(bs))
	t.Fail()
}
