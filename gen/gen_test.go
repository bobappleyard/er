package gen

import (
	. "github.com/bobappleyard/er"
	"github.com/bobappleyard/er/l2p"
	"io/ioutil"
	"os"
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
	c.Dependency = Dependency{Rel: c.Relationships[0], Sequence: true}
	d.Dependency = Dependency{Rel: d.Relationships[0]}

	generateAndTest(t, m)
}

func TestBoot(t *testing.T) {
	m := &EntityModel{
		Name: "entity_model",
		Types: []*EntityType{
			{Name: "entity_type"},
			{Name: "attribute"},
			{Name: "relationship"},
			{Name: "dependency"},
		},
	}

	entityType := m.Types[0]
	attribute := m.Types[1]
	relationship := m.Types[2]
	dependency := m.Types[3]

	entityType.Attributes = []*Attribute{
		{Name: "name", Type: StringType, Identifying: true, Owner: entityType},
	}
	attribute.Attributes = []*Attribute{
		{Name: "name", Type: StringType, Identifying: true, Owner: attribute},
		{Name: "type", Type: StringType, Owner: attribute},
		{Name: "identifying", Type: BoolType, Owner: attribute},
	}
	relationship.Attributes = []*Attribute{
		{Name: "name", Type: StringType, Identifying: true, Owner: relationship},
		{Name: "identifying", Type: BoolType, Owner: relationship},
		{Name: "path", Type: StringType, Owner: relationship},
	}
	dependency.Attributes = []*Attribute{
		{Name: "sequence", Type: BoolType, Owner: dependency},
	}

	attribute.Relationships = []*Relationship{
		{Name: "owner", Source: attribute, Target: entityType, Identifying: true},
	}
	relationship.Relationships = []*Relationship{
		{Name: "source", Source: relationship, Target: entityType, Identifying: true},
		{Name: "target", Source: relationship, Target: entityType},
	}
	dependency.Relationships = []*Relationship{
		{Name: "entity_type", Source: dependency, Target: entityType, Identifying: true},
		{Name: "relationship", Source: dependency, Target: relationship},
	}

	attribute.Dependency.Rel = attribute.Relationships[0]
	relationship.Dependency.Rel = relationship.Relationships[0]
	dependency.Dependency.Rel = dependency.Relationships[0]

	generateAndTest(t, m)
}

func generateAndTest(t *testing.T, m *EntityModel) {
	l2p.LogicalToPhysical(m)
	bs, err := generate(m)
	if err != nil {
		t.Error(err)
		return
	}
	ioutil.WriteFile(path.Join("test", m.Name, "model.go"), bs, 0777)
	cmd := exec.Command("go", "test", "./"+path.Join("test", m.Name))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

}
