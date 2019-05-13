package diagram

import (
	"fmt"
	"testing"
)

func TestPerm(t *testing.T) {
	for p := firstPerm(3); p != nil; p = nextPerm(p) {
		fmt.Println(p)
	}
}
