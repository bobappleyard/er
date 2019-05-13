package diagram

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/bobappleyard/er/l2p"
	"os"
	"testing"

	. "github.com/bobappleyard/er"
)

func TestGenerate(t *testing.T) {
	m := &EntityModel{
		Name: "square",
		Types: []*EntityType{
			{Name: "a"},
			{Name: "f"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
			{Name: "e"},
		},
	}
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
	a, b, c, d, e := m.Types[0], m.Types[2], m.Types[5], m.Types[4], m.Types[3]
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
	e.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      e,
			Target:      a,
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
	e.DependsOn = e.Relationships[0]
	f, err := os.Create("test.svg")
	if err != nil {
		t.Error(err)
		return
	}
	l2p.LogicalToPhysical(m)
	tw := towersFor(m)
	tw.calcLayout(0, 0)
	fmt.Println(tw)
	t.Fail()
	Generate(svg.New(f), m)
	f.Close()
}
