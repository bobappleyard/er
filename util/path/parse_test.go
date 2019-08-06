package path

import (
	"github.com/bobappleyard/er/util/lex"
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	var (
		tok lex.Token
		err error
	)
	inp := []byte("a/s/~b&b_name/~name")
	outp := []int{
		idTok,
		joinTok,
		idTok,
		joinTok,
		invTok,
		idTok,
		andTok,
		idTok,
		joinTok,
		invTok,
		idTok,
	}
	for _, id := range outp {
		tok, err = lexicon.Match(inp, tok)
		if err != nil {
			t.Fatal(err)
		}
		if tok.ID != id {
			t.Fatalf("expected %d, got %d", id, tok.ID)
		}
	}
	if tok.End < len(inp) {
		t.Error("failed to consume input")
	}
}

func TestParse(t *testing.T) {
	for _, test := range []struct {
		name string
		inp  string
		outp Path
	}{
		{
			name: "square",
			inp:  `~a/*/c&a/s/~b&b_name/~name`,
			outp: Intersection{
				Intersection{
					Join{
						Inverse{Term{"a"}},
						Join{Term{"*"}, Term{"c"}},
					},
					Join{Term{"a"}, Join{Term{"s"}, Inverse{Term{"b"}}}},
				},
				Join{Term{"b_name"}, Inverse{Term{"name"}}},
			},
		},
		{
			name: "arcRec",
			inp:  `owner|parent/scope`,
			outp: Union{Term{"owner"}, Join{Term{"parent"}, Term{"scope"}}},
		}, {
			name: "unionInv",
			inp:  `~(a|b)`,
			outp: Inverse{Union{Term{"a"}, Term{"b"}}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			path, err := Parse([]byte(test.inp))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(path, test.outp) {
				t.Errorf("failed: %#v != %#v", path, test.outp)
			}
		})
	}
}
