package entity_model

import (
	"testing"
)

func TestModel(t *testing.T) {
	m := New()
	err := m.Unmarshal([]byte(`
	entity_type {
		name: "a"
		attribute {
			name: "name"
			type: "string"
			identifying: "true"
		}
		relationship {
			name: "s"
			target_name: "b"
		}
	}

	entity_type {
		name: "b"
		attribute {
			name: "name"
			type: "string"
			identifying: "true"
		}
	}

	entity_type {
		name: "c"
		dependency {
			relationship_name: "parent"
		}
		attribute {
			name: "name"
			type: "string"
			identifying: "true"
		}
		relationship {
			name: "parent"
			target_name: "a"
		}
		relationship {
			name: "f"
			target_name: "d"
			constraint {
				diagonal { relationship_name: "parent" }
				diagonal { relationship_name: "s" }
				riser { relationship_name: "parent" }
			}
		}
	}

	entity_type {
		name: "d"
		dependency {
			relationship_name: "parent"
		}
		attribute {
			name: "name"
			type: "string"
			identifying: "true"
		}
		relationship {
			name: "parent"
			target_name: "b"
			identifying: "true"
		}
	}

	entity_type {
		name: "e"
		dependency {
			relationship_name: "parent"
		}
		relationship {
			name: "parent"
			target_name: "b"
		}
	}
	`))
	if err != nil {
		t.Error(err)
		return
	}
	deps := map[string]string{
		"c": "a",
		"d": "b",
		"e": "b",
	}
	m.EntityType.ForEach(func(e EntityType) error {
		e.Dependents().ForEach(func(f EntityType) error {
			if e.Name != deps[f.Name] {
				t.Errorf("got %s -> %s, expecting %s -> %s", e.Name, f.Name, e.Name, deps[f.Name])
			}
			return nil
		})
		return nil
	})
	root := []string{
		"a",
		"b",
	}
	if m.RootTypes().Count() != len(root) {
		t.Errorf("got %d root types, expecting %d", m.RootTypes().Count(), len(root))
	}
	m.RootTypes().ForEach(func(e EntityType) error {
		for _, r := range root {
			if r == e.Name {
				return nil
			}
		}
		t.Errorf("unexpected root type %s", e.Name)
		return nil
	})
}
