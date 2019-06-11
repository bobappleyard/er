package square

import (
	"testing"
)

func TestModelCRUD(t *testing.T) {
	m := New()
	assertEntries := func(s C_Set, es []C) func(*testing.T) {
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

	m.C.Upsert(C{Name: "C3", ParentName: "A1", FName: "D1"})
	m.C.Upsert(C{Name: "C3", ParentName: "A1", FName: "D2"})
	t.Run("Upsert", assertEntries(m.C, []C{
		{Name: "C1", ParentName: "A1", FName: "D1"},
		{Name: "C2", ParentName: "A1", FName: "D2"},
		{Name: "C3", ParentName: "A1", FName: "D2"},
	}))

}
