package gen

import (
	"bytes"
	"fmt"
	"github.com/bobappleyard/er"
	"go/format"
	"strings"
)

func generate(m *er.EntityModel) ([]byte, error) {
	g := &generator{
		m: m,
	}
	for _, action := range []func() error{
		g.generateHeader,
		g.generateModelDecl,
		g.generateModelCRUD,
		g.generateModelIO,
		g.generateEntities,
	} {
		if err := action(); err != nil {
			return nil, err
		}
	}
	return format.Source(g.dest.Bytes())
}

type generator struct {
	dest bytes.Buffer
	m    *er.EntityModel
}

func (g *generator) out(form string, args ...interface{}) {
	fmt.Fprintf(&g.dest, form+"\n", args...)
}

func (g *generator) generateHeader() error {
	g.out("package %s", g.m.Name)
	g.out("import (")
	g.out("%q", "github.com/bobappleyard/er")
	g.out("%q", "github.com/bobappleyard/er/rtl")
	g.out(")")
	return nil
}

func (g *generator) generateModelDecl() error {
	g.out("type Model struct {")
	for _, t := range g.m.Types {
		g.out("%s setOf%s", goName(t.Name), goName(t.Name))
	}
	g.out("}")
	g.out("func New() *Model{ ")
	g.out("m := new(Model)")
	for _, t := range g.m.Types {
		g.out("m.%s.init(m)", goName(t.Name))
	}
	g.out("return m")
	g.out("}")
	return nil
}

func (g *generator) generateModelIO() error {
	g.out("func (m *Model) Unmarshal(bs []byte) error {")
	g.out("p := rtl.NewReader(bs)")
	g.out("for p.Next() {")
	g.out("switch p.Name() {")
	for _, t := range g.dependants(nil) {
		g.out("case %q:", t.Name)
		g.out("m.%s.parse(p.Record())", goName(t.Name))
	}
	g.out("}}")
	g.out("p.ExpectEOF()")
	g.out("return p.Err()")
	g.out("}")
	return nil
}

func (g *generator) generateModelCRUD() error {
	g.out("func (m *Model) Validate() error {")
	for _, t := range g.m.Types {
		g.out("if err := m.%s.validate(); err != nil { return err }", goName(t.Name))
	}
	g.out("return nil")
	g.out("}")
	return nil
}

func (g *generator) generateEntities() error {
	for _, t := range g.m.Types {
		g.out("")
		for _, action := range []func(*er.EntityType) error{
			g.generateDecls,
			g.generateRelationships,
			g.generateCRUD,
			g.generateIO,
		} {
			if err := action(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) generateDecls(t *er.EntityType) error {
	g.out("type %s struct {", goName(t.Name))
	for _, a := range t.Attributes {
		g.out("%s %s", goName(a.Name), attrType(a))
	}
	g.out("")
	g.out("model *Model")
	g.out("}")
	g.out("")
	g.out("type attrsOf%s struct {", goName(t.Name))
	for _, a := range t.Attributes {
		g.out("%s %s", goName(a.Name), attrType(a))
	}
	g.out("}")
	g.out("")
	g.out("type setOf%s struct {", goName(t.Name))
	for _, a := range t.Attributes {
		g.out("%s %s", goName(a.Name), columnType(a))
	}
	g.out("")
	g.out("model *Model")
	g.out("query *rtl.Query")
	g.out("rows []attrsOf%s", goName(t.Name))
	g.out("}")
	g.out("func(s *setOf%s) init(m *Model) {", goName(t.Name))
	g.out("s.model = m")
	for i, a := range t.Attributes {
		init := "Column"
		if a.Identifying {
			init = "Index"
		}
		g.out("s.%s = %s%s(%d, func(idx int) %s { return s.rows[idx].%[1]s})", goName(a.Name), columnType(a), init, i, attrType(a))
	}
	g.out("}")
	g.out("")
	return nil
}

func (g *generator) generateRelationships(t *er.EntityType) error {
	g.out("func (s setOf%s) validate() error {", goName(t.Name))
	if len(t.Relationships) != 0 {
		g.out("if err := s.ForEach(func(e %s) error {", goName(t.Name))
		for _, r := range t.Relationships {
			g.out("{")
			g.out("q := e.queryFor%s()", goName(r.Name))
			g.out("if q.Count() != 1 { return er.ErrMissingEntity }")
			if len(r.Constraints) == 0 {
				g.out("}")
				continue
			}
			g.out("t := q.ExactlyOne()")
			for _, c := range r.Constraints {
				diagonal := make([]string, len(c.Diagonal.Components))
				for i, m := range c.Diagonal.Components {
					diagonal[i] = goName(m.Rel.Name) + "()"
				}
				riser := make([]string, len(c.Riser.Components))
				for i, m := range c.Riser.Components {
					riser[i] = goName(m.Rel.Name) + "()"
				}
				g.out("if e.%s != t.%s {", strings.Join(diagonal, "."), strings.Join(riser, "."))
				g.out("return er.ErrMissingEntity")
				g.out("}")
			}
			g.out("}")
		}
		g.out("return nil")
		g.out("}); err != nil { return err }")
	}
	g.out("return nil")
	g.out("}")
	g.out("")
	for _, r := range t.Relationships {
		g.out("func (e %s) %s() %s {", goName(t.Name), goName(r.Name), goName(r.Target.Name))
		g.out("return e.queryFor%s().ExactlyOne()", goName(r.Name))
		g.out("}")
		g.out("")

		g.out("func (e %s) queryFor%s() setOf%s {", goName(t.Name), goName(r.Name), goName(r.Target.Name))
		g.out("var q rtl.Query")
		for _, k := range r.Implementation {
			path := make([]string, len(k.BasePath)+1)
			for i, c := range k.BasePath {
				path[i] = goName(c.Rel.Name) + "()"
			}
			path[len(path)-1] = goName(k.Source.Name)
			g.out("q = q.And(e.model.%s.%s.Eq(e.%s))", goName(r.Target.Name), goName(k.Target.Name), strings.Join(path, "."))
		}
		g.out("return e.model.%s.Where(q)", goName(r.Target.Name))
		g.out("}")
		g.out("")
	}
	return nil
}

func (g *generator) generateCRUD(t *er.EntityType) error {
	g.out("func (s setOf%s) ForEach(f func(%[1]s) error) error {", goName(t.Name))
	g.out("q := rtl.All(len(s.rows))")
	g.out("if s.query != nil { q = rtl.EvalQuery(*s.query, len(s.rows)) }")
	g.out("for q.Next() {")
	g.out("d := s.rows[q.This()]")
	g.out("if err := f(%s{", goName(t.Name))
	g.out("model: s.model,")
	for _, a := range t.Attributes {
		g.out("%s: d.%[1]s,", goName(a.Name))
	}
	g.out("}); err != nil { return err }")
	g.out("}")
	g.out("return nil")
	g.out("}")
	g.out("")

	g.out("func (s setOf%s) Count() int {", goName(t.Name))
	g.out("c := 0")
	g.out("s.ForEach(func(%s) error {", goName(t.Name))
	g.out("c++")
	g.out("return nil")
	g.out("})")
	g.out("return c")
	g.out("}")

	g.out("func (s setOf%s) ExactlyOne() %[1]s {", goName(t.Name))
	g.out("var res %s", goName(t.Name))
	g.out("s.ForEach(func(t %s) error {", goName(t.Name))
	g.out("res = t")
	g.out("return nil")
	g.out("})")
	g.out("return res")
	g.out("}")

	g.out("func (s setOf%s) Where(q rtl.Query) setOf%[1]s {", goName(t.Name))
	g.out("res := s")
	g.out("if res.query != nil { q = q.And(*res.query) }")
	g.out("res.query = &q")
	g.out("return res")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) Insert(e %[1]s) error {", goName(t.Name))
	g.out("if s.query != nil { return er.ErrImmutableSet }")
	g.out("r := s.evalKey(e)")
	g.out("if r.Next() { return er.ErrDuplicateKey }")
	g.out("s.clearSpace(r)")
	g.out("s.writeRow(r, e)")
	g.out("return nil")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) Update(e %[1]s) error {", goName(t.Name))
	g.out("if s.query != nil { return er.ErrImmutableSet }")
	g.out("r := s.evalKey(e)")
	g.out("if !r.Next() { return er.ErrMissingEntity }")
	g.out("s.writeRow(r, e)")
	g.out("return nil")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) Upsert(e %[1]s) error {", goName(t.Name))
	g.out("if s.query != nil { return er.ErrImmutableSet }")
	g.out("r := s.evalKey(e)")
	g.out("if !r.Next() { s.clearSpace(r) }")
	g.out("s.writeRow(r, e)")
	g.out("return nil")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) Delete(e %[1]s) error {", goName(t.Name))
	g.out("if s.query != nil { return er.ErrImmutableSet }")
	g.out("r := s.evalKey(e)")
	g.out("if !r.Next() { return er.ErrMissingEntity }")
	g.out("copy(s.rows[r.This():], s.rows[r.This():+1])")
	g.out("s.rows = s.rows[:len(s.rows)-1]")
	g.out("return nil")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) evalKey(e %[1]s) *rtl.QueryResult {", goName(t.Name))
	g.out("var query rtl.Query")
	for _, a := range t.Attributes {
		if !a.Identifying {
			continue
		}
		g.out("query = query.And(s.%s.Eq(e.%[1]s))", goName(a.Name))
	}
	g.out("return rtl.EvalQuery(query, len(s.rows))")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) clearSpace(r *rtl.QueryResult) {", goName(t.Name))
	g.out("s.rows = append(s.rows, attrsOf%s{})", goName(t.Name))
	g.out("copy(s.rows[r.This()+1:], s.rows[r.This():])")
	g.out("}")
	g.out("")

	g.out("func (s *setOf%s) writeRow(r *rtl.QueryResult, e %[1]s) {", goName(t.Name))
	g.out("s.rows[r.This()] = attrsOf%s {", goName(t.Name))
	for _, a := range t.Attributes {
		g.out("%s: e.%[1]s,", goName(a.Name))
	}
	g.out("}")
	g.out("}")
	g.out("")

	return nil
}

func (g *generator) generateIO(t *er.EntityType) error {
	omit := map[string]bool{}
	if t.DependsOn != nil {
		parent := t.DependsOn.Target
		g.out("func (s *setOf%s) parse(p *rtl.Reader, parent %s) {", goName(t.Name), goName(parent.Name))
		g.out("var e %s", goName(t.Name))
		for _, k := range t.DependsOn.Implementation {
			g.out("e.%s = parent.%s", goName(k.Source.Name), goName(k.Target.Name))
			omit[k.Source.Name] = true
		}
	} else {
		g.out("func (s *setOf%s) parse(p *rtl.Reader) {", goName(t.Name))
		g.out("var e %s", goName(t.Name))
	}
	g.out("for p.Next() { switch p.Name() {")
	for _, a := range t.Attributes {
		if omit[a.Name] {
			continue
		}
		g.out("case %q: e.%s = p.%s()", a.Name, goName(a.Name), attrParse(a))
	}
	for _, d := range g.dependants(t) {
		g.out("case %q: s.model.%s.parse(p.Record(), e)", d.Name, goName(d.Name))
	}
	g.out("}}")
	g.out("if p.Err() == nil { p.SetErr(s.Insert(e)) }")
	g.out("}")
	return nil
}

func (g *generator) dependants(t *er.EntityType) []*er.EntityType {
	var res []*er.EntityType
	for _, u := range g.m.Types {
		var p *er.EntityType
		if u.DependsOn != nil {
			p = u.DependsOn.Target
		}
		if t == p {
			res = append(res, u)
		}
	}
	return res
}

func goName(name string) string {
	parts := strings.Split(name, "_")
	for i, p := range parts {
		parts[i] = strings.ToTitle(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func attrType(a *er.Attribute) string {
	switch a.Type {
	case er.IntType:
		return "int"
	case er.FloatType:
		return "float64"
	case er.StringType:
		return "string"
	}
	return "?"
}

func columnType(a *er.Attribute) string {
	var tn string
	switch a.Type {
	case er.IntType:
		tn = "rtl.Int"
	case er.FloatType:
		tn = "rtl.Float64"
	case er.StringType:
		tn = "rtl.String"
	default:
		return "?"
	}
	return tn
}

func attrParse(a *er.Attribute) string {
	var tn string
	switch a.Type {
	case er.IntType:
		tn = "IntAttr"
	case er.FloatType:
		tn = "FloatAttr"
	case er.StringType:
		tn = "StringAttr"
	default:
		return "?"
	}
	return tn
}
