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
						InverseTerm{"a"},
						Join{Term{"*"}, Term{"c"}},
					},
					Join{Term{"a"}, Join{Term{"s"}, InverseTerm{"b"}}},
				},
				Join{Term{"b_name"}, InverseTerm{"name"}},
			},
		},
		{
			name: "arcRec",
			inp:  `owner|parent/scope`,
			outp: Union{Term{"owner"}, Join{Term{"parent"}, Term{"scope"}}},
		},
		{
			name: "squareInv",
			inp:  `~(a/s/~b&b_name/~name)`,
			outp: Intersection{
				Join{Join{Term{"b"}, InverseTerm{"s"}}, InverseTerm{"a"}},
				Join{Term{"name"}, InverseTerm{"b_name"}},
			},
		}, {
			name: "unionInv",
			inp:  `~(a|b)`,
			outp: Union{InverseTerm{"a"}, InverseTerm{"b"}},
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
