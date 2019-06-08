package square

import (
	"testing"
)

func TestModelCRUD(t *testing.T) {
	m := New()
	m.A.Insert(A{Name: "A1", SName: "B1"})
	m.B.Insert(B{Name: "B1"})
	m.C.Insert(C{Name: "C2", ParentName: "A1", FName: "D1"})
	m.C.Insert(C{Name: "C1", ParentName: "A1", FName: "D1"})
	m.D.Insert(D{Name: "D1", ParentName: "B1"})
	m.Validate()
	m.C.ForEach(func(c C) error {
		t.Error(c)
		return nil
	})
	t.Error(m)
}
