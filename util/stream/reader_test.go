package stream

import (
	"fmt"
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
			name:    "NoName",
			in:      "{",
			success: false,
		},
		{
			name:    "IllegalChar",
			in:      "@",
			success: false,
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
	p := NewReader([]byte(`"value"`))
	p.kind = Field
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

func TestCodegen(t *testing.T) {
	type member struct {
		empNo    int
		name     string
		goldStar bool
	}

	type team struct {
		name        string
		leaderEmpNo int
		members     []member
	}

	createMember := func(t *team) *member {
		t.members = append(t.members, member{})
		return &t.members[len(t.members)-1]
	}

	parseMember := func(m *member, r *Reader) {
		for r.Next() {
			switch r.Name() {
			case "emp_no":
				m.empNo = r.IntField()
			case "name":
				m.name = r.StringField()
			case "gold_star":
				m.goldStar = r.BoolField()
			}
		}
	}

	parseTeam := func(t *team, r *Reader) {
		for r.Next() {
			switch r.Name() {
			case "name":
				t.name = r.StringField()
			case "leader_emp_no":
				t.leaderEmpNo = r.IntField()
			case "member":
				m := createMember(t)
				parseMember(m, r.Record())
			}
		}
	}

	r := NewReader([]byte(`
team {
	name: "technology"
	leader_emp_no: "1"

	# this is a comment, and so ignored
	member {
		emp_no: "1"
		name: "big boss"
		gold_star: "true"
	}
}
	`))
	var ts []team
	for r.Next() {
		switch r.Name() {
		case "team":
			var t team
			parseTeam(&t, r.Record())
			ts = append(ts, t)
		}
	}
	r.ExpectEOF()
	if r.Err() != nil {
		t.Fatalf("error was: %s", r.Err())
	}
	if ts[0].members[0].name != "big boss" {
		t.Errorf("bad name")
	}
}

func ExampleReader_Next() {
	r := NewReader([]byte(`a{} a{b{c: "value"}}`))
	var parse func(r *Reader)
	parse = func(r *Reader) {
		for r.Next() {
			switch r.Kind() {
			case Field:
				fmt.Println(r.Name(), r.StringField())
			case Record:
				name := r.Name()
				fmt.Println("record", name)
				parse(r.Record())
				fmt.Println("end", name)
			}
		}
	}
	parse(r)
	// Output:
	// record a
	// end a
	// record a
	// record b
	// c value
	// end b
	// end a
}
