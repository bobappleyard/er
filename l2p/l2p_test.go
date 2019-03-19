package er

import (
	"reflect"
	"sort"
	"testing"

	. "github.com/bobappleyard/er"
)

func TestLine(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
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
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	b.Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: b,
			Target: a,
		},
		{
			Name:   "f",
			Source: b,
			Target: c,
		},
	}
	c.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      c,
			Target:      a,
			Identifying: true,
		},
	}
	logicalToPhysical(&m)
	testAttrs(t, m.Types[1].Attributes, []string{"name", "parent_name", "f_name", "f_parent_name"})
}

func TestTriangle(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
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
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	b.Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: b,
			Target: a,
		},
		{
			Name:   "f",
			Source: b,
			Target: c,
		},
	}
	c.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      c,
			Target:      a,
			Identifying: true,
		},
	}
	b.Relationships[1].Constraints = []Constraint{
		{
			Diagonal{Components: []Component{
				{Rel: b.Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: c.Relationships[0]},
			}},
		},
	}

	logicalToPhysical(&m)
	testAttrs(t, m.Types[1].Attributes, []string{"name", "parent_name", "f_name"})
}

func TestSquare(t *testing.T) {
	m := EntityModel{
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

	logicalToPhysical(&m)
	testAttrs(t, m.Types[2].Attributes, []string{"name", "parent_name", "f_name"})
}

func TestSquareLongRiser(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
			{Name: "e"},
		},
	}
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	d := m.Types[3]
	e := m.Types[4]
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
			Target: e,
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
			Target:      d,
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
				{Rel: e.Relationships[0]},
				{Rel: d.Relationships[0]},
			}},
		},
	}

	logicalToPhysical(&m)
	testAttrs(t, m.Types[2].Attributes, []string{"name", "parent_name", "f_name", "f_parent_name"})
}

func TestCube(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
			{Name: "e"},
			{Name: "f"},
		},
	}
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	d := m.Types[3]
	e := m.Types[4]
	f := m.Types[5]
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
		{
			Name:   "parent2",
			Source: c,
			Target: e,
		},
	}
	d.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      d,
			Target:      b,
			Identifying: true,
		},
		{
			Name:        "parent2",
			Source:      d,
			Target:      f,
			Identifying: true,
		},
	}
	e.Relationships = []*Relationship{
		{
			Name:   "s2",
			Source: e,
			Target: f,
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
		{
			Diagonal{[]Component{
				{Rel: c.Relationships[2]},
				{Rel: e.Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: d.Relationships[1]},
			}},
		},
	}

	logicalToPhysical(&m)
	testAttrs(t, m.Types[2].Attributes, []string{"name", "f_name", "parent_name", "parent2_name"})
}

func TestCubeShared(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
			{Name: "e"},
		},
	}
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
	d := m.Types[3]
	e := m.Types[4]
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
		{
			Name:   "s2",
			Source: a,
			Target: e,
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
		{
			Name:        "parent2",
			Source:      d,
			Target:      e,
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
		{
			Diagonal{[]Component{
				{Rel: c.Relationships[0]},
				{Rel: a.Relationships[1]},
			}},
			Riser{[]Component{
				{Rel: d.Relationships[1]},
			}},
		},
	}

	logicalToPhysical(&m)
	testAttrs(t, m.Types[2].Attributes, []string{"name", "f_name", "parent_name"})
}

func TestTriangleTwin(t *testing.T) {
	m := EntityModel{
		Types: []*EntityType{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
	}
	a := m.Types[0]
	b := m.Types[1]
	c := m.Types[2]
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
	b.Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: b,
			Target: a,
		},
		{
			Name:   "f",
			Source: b,
			Target: c,
		},
	}
	c.Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      c,
			Target:      a,
			Identifying: true,
		},
		{
			Name:        "parent2",
			Source:      c,
			Target:      a,
			Identifying: true,
		},
	}
	b.Relationships[1].Constraints = []Constraint{
		{
			Diagonal{Components: []Component{
				{Rel: b.Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: c.Relationships[0]},
			}},
		},
		{
			Diagonal{Components: []Component{
				{Rel: b.Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: c.Relationships[1]},
			}},
		},
	}

	logicalToPhysical(&m)
	testAttrs(t, m.Types[1].Attributes, []string{"name", "parent_name", "f_name"})
}

func testAttrs(t *testing.T, attrs []*Attribute, names []string) {
	var anames []string
	for _, a := range attrs {
		anames = append(anames, a.Name)
	}
	sort.Strings(names)
	sort.Strings(anames)
	if !reflect.DeepEqual(names, anames) {
		t.Errorf("expected attrs: %s, got attrs: %s", names, anames)
	}
}
