package square

import (
	"testing"
)

func TestModelCRUD(t *testing.T) {
	m := Model{}
	m.A.Insert(A{Name: "A1", SName: "B1"})
	m.B.Insert(B{Name: "B1"})
	m.C.Insert(C{Name: "C2", ParentName: "A1", FName: "D1"})
	m.C.Insert(C{Name: "C1", ParentName: "A1", FName: "D1"})
	m.D.Insert(D{Name: "D1", ParentName: "B1"})
	m.Validate()
	t.Error(m)
}
