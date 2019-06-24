package rtl

import (
	"testing"
)

func TestReaderNext(t *testing.T) {
	for _, test := range []struct {
		name, in, out string
		success       bool
	}{
		{
			name:    "Empty",
			in:      "",
			success: false,
		},
		{
			name:    "GoodInput",
			in:      "name ",
			out:     "name",
			success: true,
		},
		{
			name:    "Attr",
			in:      "name:",
			out:     "name",
			success: true,
		},
		{
			name:    "Rec",
			in:      "name{",
			out:     "name",
			success: true,
		},
		{
			name:    "LeadingSpace",
			in:      "  name{",
			out:     "name",
			success: true,
		},
		{
			name:    "NonAlnum",
			in:      "pH7:",
			out:     "pH7",
			success: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := NewReader([]byte(test.in))
			success := p.Next()
			if success != test.success {
				t.Errorf("success was %t, expecting %t", success, test.success)
			}
			if p.name != test.out {
				t.Errorf("name was %q, expecting %q", p.name, test.out)
			}
		})
	}
}

func TestParseAttr(t *testing.T) {
	p := NewReader([]byte(`:"value"`))
	p.attrsOK = true
	value := p.parseAttr()
	if value != `"value"` {
		t.Log(value)
		t.Error("StringAttr() failed")
	}
	if p.Err() != nil {
		t.Log(p.Err())
		t.Error("Err() failed")
	}
}
