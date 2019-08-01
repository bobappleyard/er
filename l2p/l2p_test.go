package l2p

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
	LogicalToPhysical(&m)
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
			Path:   "parent/~parent",
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

	LogicalToPhysical(&m)
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
			Path:   "parent/s/~parent",
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

	LogicalToPhysical(&m)
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
			Path:   "parent/s/~parent/~parent",
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

	LogicalToPhysical(&m)
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
			Path:   "parent/s/~parent&parent2/s2/~parent2",
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

	LogicalToPhysical(&m)
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
			Path:   "parent/(s/~parent&s2/~parent2)",
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

	LogicalToPhysical(&m)
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
			Path:   "parent/(~parent&~parent2)",
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

	LogicalToPhysical(&m)
	testAttrs(t, m.Types[1].Attributes, []string{"name", "parent_name", "f_name"})
}

// func TestSequence(t *testing.T) {
// 	m := EntityModel{
// 		Types: []*EntityType{
// 			{Name: "a"},
// 			{Name: "b"},
// 			{Name: "c"},
// 		},
// 	}
// 	a := m.Types[0]
// 	b := m.Types[1]
// 	c := m.Types[2]
// 	for _, t := range m.Types {
// 		t.Attributes = []*Attribute{
// 			{
// 				Owner:       t,
// 				Name:        "name",
// 				Type:        StringType,
// 				Identifying: true,
// 			},
// 		}
// 	}
// 	b.Relationships = []*Relationship{
// 		{
// 			Name:   "parent",
// 			Source: b,
// 			Target: a,
// 		},
// 		{
// 			Name:   "f",
// 			Source: b,
// 			Target: c,
// 		},
// 	}
// 	c.Relationships = []*Relationship{
// 		{
// 			Name:        "parent",
// 			Source:      c,
// 			Target:      a,
// 			Identifying: true,
// 		},
// 	}
// 	b.Relationships[1].Constraints = []Constraint{
// 		{
// 			Diagonal{Components: []Component{
// 				{Rel: b.Relationships[0]},
// 			}},
// 			Riser{[]Component{
// 				{Rel: c.Relationships[0]},
// 			}},
// 		},
// 	}
// 	b.Dependency = Dependency{
// 		Rel:      b.Relationships[0],
// 		Sequence: true,
// 	}

// 	LogicalToPhysical(&m)
// 	testAttrs(t, m.Types[1].Attributes, []string{"name", "parent_name", "f_name", "seq"})
// }

// func TestSequenceKey(t *testing.T) {
// 	m := EntityModel{
// 		Types: []*EntityType{
// 			{Name: "a"},
// 			{Name: "b"},
// 			{Name: "c"},
// 		},
// 	}
// 	a := m.Types[0]
// 	b := m.Types[1]
// 	c := m.Types[2]
// 	for _, t := range m.Types {
// 		t.Attributes = []*Attribute{
// 			{
// 				Owner:       t,
// 				Name:        "name",
// 				Type:        StringType,
// 				Identifying: true,
// 			},
// 		}
// 	}
// 	b.Relationships = []*Relationship{
// 		{
// 			Name:   "parent",
// 			Source: b,
// 			Target: a,
// 		},
// 		{
// 			Name:   "f",
// 			Source: b,
// 			Target: c,
// 		},
// 	}
// 	c.Relationships = []*Relationship{
// 		{
// 			Name:        "parent",
// 			Source:      c,
// 			Target:      a,
// 			Identifying: true,
// 		},
// 	}
// 	b.Relationships[1].Constraints = []Constraint{
// 		{
// 			Diagonal{Components: []Component{
// 				{Rel: b.Relationships[0]},
// 			}},
// 			Riser{[]Component{
// 				{Rel: c.Relationships[0]},
// 			}},
// 		},
// 	}
// 	c.Dependency = Dependency{
// 		Rel:      c.Relationships[0],
// 		Sequence: true,
// 	}

// 	LogicalToPhysical(&m)
// 	testAttrs(t, m.Types[1].Attributes, []string{"f_name", "f_seq", "name", "parent_name"})
// }

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
