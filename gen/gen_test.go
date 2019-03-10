package gen

import (
	"bytes"
	"github.com/bobappleyard/er"
	"github.com/bobappleyard/er/l2p"
	"go/format"
	"go/token"
	"testing"
)

func TestGen(t *testing.T) {
	m := &er.EntityModel{
		Name: "square",
		Types: []*er.EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
		},
	}
	for _, t := range m.Types {
		t.Attributes = []*er.Attribute{
			{
				Owner:       t,
				Name:        "name",
				Type:        er.StringType,
				Identifying: true,
			},
		}
	}
	m.Types[0].Relationships = []*er.Relationship{
		{
			Name:   "s",
			Source: m.Types[0],
			Target: m.Types[1],
		},
	}
	m.Types[2].Relationships = []*er.Relationship{
		{
			Name:   "parent",
			Source: m.Types[2],
			Target: m.Types[0],
		},
		{
			Name:   "f",
			Source: m.Types[2],
			Target: m.Types[3],
		},
	}
	m.Types[3].Relationships = []*er.Relationship{
		{
			Name:        "parent",
			Source:      m.Types[3],
			Target:      m.Types[1],
			Identifying: true,
		},
	}
	m.Types[2].Relationships[1].Constraints = []er.Constraint{
		{
			er.Diagonal{[]er.Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[0]},
			}},
			er.Riser{[]er.Component{
				{Rel: m.Types[3].Relationships[0]},
			}},
		},
	}
	l2p.LogicalToPhysical(m)
	p, err := generate(m)
	if err != nil {
		t.Error(err)
		return
	}
	var b bytes.Buffer
	b.WriteByte('\n')
	format.Node(&b, token.NewFileSet(), p)
	t.Log(b.String())
	t.Fail()
}
