package diagram

import (
	"fmt"
)

func (t *tower) Format(state fmt.State, c rune) {
	prefix := ""
	for c := t.up; c != nil; c = c.up {
		prefix += "  "
	}
	printf := func(format string, args ...interface{}) {
		fmt.Fprintf(state, prefix+format+"\n", args...)
	}
	nm := ""
	if t.t != nil {
		nm = t.t.Name + " "
	}
	printf("%s{", nm)
	printf("  head: %v", t.head)
	printf("  body: %v", t.body)
	for _, t := range t.down {
		t.Format(state, c)
	}
	printf("}")
}
