package square

import (
	"testing"
)

type cIter interface {
	ForEach(func(C) error) error
}

func TestModelCRUD(t *testing.T) {

	assertRels := func(s cIter, fs []D) func(*testing.T) {
		return func(t *testing.T) {
			i := 0
			s.ForEach(func(c C) error {
				if i < len(fs) {
					d := c.F()
					if d.Name != fs[i].Name ||
						d.ParentName != fs[i].ParentName {
						t.Errorf("[%d] got %v, expecting %v", i, d, fs[i])
					}
				}
				i++
				return nil
			})
			if i != len(fs) {
				t.Errorf("got %d entities, expecting %d", i, len(fs))
			}
		}
	}

	m := New()

	m.A.Insert(A{Name: "A1", SName: "B1"})
	m.B.Insert(B{Name: "B1"})
	m.C.Insert(C{Name: "C2", ParentName: "A1", FName: "D1"})
	m.C.Insert(C{Name: "C1", ParentName: "A1", FName: "D1"})
	m.D.Insert(D{Name: "D1", ParentName: "B1"})
	t.Run("Insert", assertEntries(m.C, []C{
		{Name: "C1", ParentName: "A1", FName: "D1"},
		{Name: "C2", ParentName: "A1", FName: "D1"},
	}))

	m.D.Insert(D{Name: "D2", ParentName: "B1"})
	m.C.Update(C{Name: "C2", ParentName: "A1", FName: "D2"})
	t.Run("Update", assertEntries(m.C, []C{
		{Name: "C1", ParentName: "A1", FName: "D1"},
		{Name: "C2", ParentName: "A1", FName: "D2"},
	}))
	t.Run("Select1", assertEntries(m.C.Where(m.C.ParentName.Eq("A1")), []C{
		{Name: "C1", ParentName: "A1", FName: "D1"},
		{Name: "C2", ParentName: "A1", FName: "D2"},
	}))
	t.Run("Select2", assertEntries(m.C.Where(m.C.ParentName.Eq("A1").And(m.C.FName.Eq("D2"))), []C{
		{Name: "C2", ParentName: "A1", FName: "D2"},
	}))

	t.Run("Rels", assertRels(m.C, []D{
		{Name: "D1", ParentName: "B1"},
		{Name: "D2", ParentName: "B1"},
	}))

	if err := m.Validate(); err != nil {
		t.Error("validation failed on valid model")
	}
	m.C.Insert(C{Name: "C4", ParentName: "A1", FName: "D3"})
	if err := m.Validate(); err == nil {
		t.Error("validation succeeded on invalid model")
	}

}

func TestModelParse(t *testing.T) {
	m := New()
	err := m.Unmarshal([]byte(`
	a {
		name: "A1"
		s_name: "B1"

		c {
			name: "C1"
			f_name: "D1"
		}
	}

	b {
		name: "B1"

		d {
			name: "D1"
		}
	}
	`))
	if err != nil {
		t.Errorf("parse failed: %v", err)
	}
	err = m.Validate()
	if err != nil {
		t.Errorf("validation failed: %v", err)
	}
	t.Run("Check", assertEntries(m.C, []C{{Name: "C1", ParentName: "A1", FName: "D1"}}))
}

func assertEntries(s cIter, es []C) func(*testing.T) {
	return func(t *testing.T) {
		i := 0
		s.ForEach(func(c C) error {
			if i < len(es) {
				if c.Name != es[i].Name ||
					c.ParentName != es[i].ParentName ||
					c.FName != es[i].FName {
					t.Errorf("[%d] got %v, expecting %v", i, c, es[i])
				}
			}
			i++
			return nil
		})
		if i != len(es) {
			t.Errorf("got %d entities, expecting %d", i, len(es))
		}
	}
}
