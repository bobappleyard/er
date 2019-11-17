package stream

import (
	"testing"
)

func TestDOM(t *testing.T) {
	var d DOM
	err := d.Unmarshal([]byte(`
	
	root {
		child {
			name: "Test"
		}
		child {
			name: "Test2"
		}
	}
	
	`))

	if err != nil {
		t.Fatal(err)
	}
	if d.Children[0].Children[0].Name != "child" {
		t.Error(d)
	}
	if d.Children[0].Children[1].Fields["name"] != "Test2" {
		t.Error(d)
	}

	buf, err := d.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != `root{child{name:"Test"}child{name:"Test2"}}` {
		t.Error(buf)
	}
}
