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
	m.Types[1].Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: m.Types[1],
			Target: m.Types[0],
		},
		{
			Name:   "f",
			Source: m.Types[1],
			Target: m.Types[2],
		},
	}
	m.Types[2].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[2],
			Target:      m.Types[0],
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
	m.Types[1].Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: m.Types[1],
			Target: m.Types[0],
		},
		{
			Name:   "f",
			Source: m.Types[1],
			Target: m.Types[2],
		},
	}
	m.Types[2].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[2],
			Target:      m.Types[0],
			Identifying: true,
		},
	}
	m.Types[1].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{Components: []Component{
				{Rel: m.Types[1].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[2].Relationships[0]},
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
	m.Types[0].Relationships = []*Relationship{
		{
			Name:   "s",
			Source: m.Types[0],
			Target: m.Types[1],
		},
	}
	m.Types[2].Relationships = []*Relationship{
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
	m.Types[3].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[3],
			Target:      m.Types[1],
			Identifying: true,
		},
	}
	m.Types[2].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[3].Relationships[0]},
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
	m.Types[0].Relationships = []*Relationship{
		{
			Name:   "s",
			Source: m.Types[0],
			Target: m.Types[1],
		},
	}
	m.Types[2].Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: m.Types[2],
			Target: m.Types[0],
		},
		{
			Name:   "f",
			Source: m.Types[2],
			Target: m.Types[4],
		},
	}
	m.Types[3].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[3],
			Target:      m.Types[1],
			Identifying: true,
		},
	}
	m.Types[4].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[4],
			Target:      m.Types[3],
			Identifying: true,
		},
	}
	m.Types[2].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[4].Relationships[0]},
				{Rel: m.Types[3].Relationships[0]},
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
	m.Types[0].Relationships = []*Relationship{
		{
			Name:   "s",
			Source: m.Types[0],
			Target: m.Types[1],
		},
	}
	m.Types[2].Relationships = []*Relationship{
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
		{
			Name:   "parent2",
			Source: m.Types[2],
			Target: m.Types[4],
		},
	}
	m.Types[3].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[3],
			Target:      m.Types[1],
			Identifying: true,
		},
		{
			Name:        "parent2",
			Source:      m.Types[3],
			Target:      m.Types[5],
			Identifying: true,
		},
	}
	m.Types[4].Relationships = []*Relationship{
		{
			Name:   "s2",
			Source: m.Types[4],
			Target: m.Types[5],
		},
	}
	m.Types[2].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[3].Relationships[0]},
			}},
		},
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[2]},
				{Rel: m.Types[4].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[3].Relationships[1]},
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
	m.Types[0].Relationships = []*Relationship{
		{
			Name:   "s",
			Source: m.Types[0],
			Target: m.Types[1],
		},
		{
			Name:   "s2",
			Source: m.Types[0],
			Target: m.Types[4],
		},
	}
	m.Types[2].Relationships = []*Relationship{
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
	m.Types[3].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[3],
			Target:      m.Types[1],
			Identifying: true,
		},
		{
			Name:        "parent2",
			Source:      m.Types[3],
			Target:      m.Types[4],
			Identifying: true,
		},
	}
	m.Types[2].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[3].Relationships[0]},
			}},
		},
		{
			Diagonal{[]Component{
				{Rel: m.Types[2].Relationships[0]},
				{Rel: m.Types[0].Relationships[1]},
			}},
			Riser{[]Component{
				{Rel: m.Types[3].Relationships[1]},
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
	m.Types[1].Relationships = []*Relationship{
		{
			Name:   "parent",
			Source: m.Types[1],
			Target: m.Types[0],
		},
		{
			Name:   "f",
			Source: m.Types[1],
			Target: m.Types[2],
		},
	}
	m.Types[2].Relationships = []*Relationship{
		{
			Name:        "parent",
			Source:      m.Types[2],
			Target:      m.Types[0],
			Identifying: true,
		},
		{
			Name:        "parent2",
			Source:      m.Types[2],
			Target:      m.Types[0],
			Identifying: true,
		},
	}
	m.Types[1].Relationships[1].Constraints = []Constraint{
		{
			Diagonal{Components: []Component{
				{Rel: m.Types[1].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[2].Relationships[0]},
			}},
		},
		{
			Diagonal{Components: []Component{
				{Rel: m.Types[1].Relationships[0]},
			}},
			Riser{[]Component{
				{Rel: m.Types[2].Relationships[1]},
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
