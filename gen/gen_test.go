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
			{Name: "implementation"},
			{Name: "constraint"},
			{Name: "step"},
			{Name: "diagonal"},
			{Name: "riser"},
		},
	}

	entityType := m.Types[0]
	attribute := m.Types[1]
	relationship := m.Types[2]
	dependency := m.Types[3]
	implementation := m.Types[4]
	constraint := m.Types[5]
	step := m.Types[6]
	diagonal := m.Types[7]
	riser := m.Types[8]

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
	}
	dependency.Attributes = []*Attribute{
		{Name: "sequence", Type: BoolType, Owner: dependency},
	}
	step.Attributes = []*Attribute{
		{Name: "relationship_name", Type: StringType, Owner: step},
	}
	diagonal.Attributes = []*Attribute{
		{Name: "relationship_name", Type: StringType, Owner: diagonal},
	}
	riser.Attributes = []*Attribute{
		{Name: "relationship_name", Type: StringType, Owner: riser},
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
	implementation.Relationships = []*Relationship{
		{Name: "relationship", Source: implementation, Target: relationship, Identifying: true},
		{Name: "target", Source: implementation, Target: attribute, Identifying: true},
		{Name: "source", Source: implementation, Target: attribute},
	}
	constraint.Relationships = []*Relationship{
		{Name: "relationship", Source: constraint, Target: relationship, Identifying: true},
	}
	step.Relationships = []*Relationship{
		{Name: "implementation", Source: step, Target: implementation, Identifying: true},
	}
	diagonal.Relationships = []*Relationship{
		{Name: "constraint", Source: diagonal, Target: constraint, Identifying: true},
	}
	riser.Relationships = []*Relationship{
		{Name: "constraint", Source: riser, Target: constraint, Identifying: true},
	}

	attribute.Dependency.Rel = attribute.Relationships[0]
	relationship.Dependency.Rel = relationship.Relationships[0]
	dependency.Dependency.Rel = dependency.Relationships[0]
	implementation.Dependency.Rel = implementation.Relationships[0]
	constraint.Dependency.Rel = constraint.Relationships[0]
	step.Dependency.Rel = step.Relationships[0]
	diagonal.Dependency.Rel = diagonal.Relationships[0]
	riser.Dependency.Rel = riser.Relationships[0]

	implementation.Dependency.Sequence = true
	constraint.Dependency.Sequence = true
	step.Dependency.Sequence = true
	diagonal.Dependency.Sequence = true
	riser.Dependency.Sequence = true

	followPath := func(e *EntityType, path ...string) []Component {
		var res []Component
		for _, step := range path {
			for _, rel := range e.Relationships {
				if rel.Name != step {
					continue
				}
				res = append(res, Component{Rel: rel})
				e = rel.Target
				break
			}
		}
		return res
	}

	dependency.Relationships[1].Constraints = []Constraint{{
		Diagonal: Diagonal{Components: followPath(dependency, "entity_type")},
		Riser:    Riser{Components: followPath(relationship, "source")},
	}}

	implementation.Relationships[1].Constraints = []Constraint{{
		Diagonal: Diagonal{Components: followPath(implementation, "relationship", "target")},
		Riser:    Riser{Components: followPath(attribute, "owner")},
	}}

	implementation.Relationships[2].Constraints = []Constraint{{
		Diagonal: Diagonal{Components: followPath(implementation, "relationship", "source")},
		Riser:    Riser{Components: followPath(attribute, "owner")},
	}}

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
