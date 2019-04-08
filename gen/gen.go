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
	g.out("%q", "github.com/bobappleyard/rsf")
	g.out(")")
	return nil
}

func (g *generator) generateModelDecl() error {
	g.out("type Model struct {")
	for _, t := range g.m.Types {
		g.out("%s %s_Set", goName(t.Name), goName(t.Name))
	}
	g.out("}")
	return nil
}

func (g *generator) generateModelIO() error {
	g.out("func (m *Model) Unmarshal(bs []byte) error {")
	g.out("p := rsf.NewReader(bs)")
	g.out("for p.Next() == rsf.RecordStart {")
	g.out("switch p.Name {")
	for _, t := range g.dependants(nil) {
		g.out("case %q:", t.Name)
		g.out("if err := m.%s.parse(p); err != nil {return err}", goName(t.Name))
	}
	g.out("default:return er.UnknownEntityType")
	g.out("}")
	g.out("if p.Next() != rsf.RecordEnd {return er.SyntaxError}")
	g.out("}")
	g.out("if p.Next() != rsf.EOF {return er.SyntaxError}")
	g.out("return nil")
	g.out("}")
	return nil
}

func (g *generator) generateEntities() error {
	for _, t := range g.m.Types {
		g.out("")
		for _, action := range []func(*er.EntityType) error{
			g.generateDecls,
			g.generateRelations,
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
	g.out("model *Model")
	for _, a := range t.Attributes {
		g.out("%s %s", goName(a.Name), attrType(a.Type))
	}
	g.out("}")
	g.out("type %s_Set struct {", goName(t.Name))
	g.out("model *Model")
	for _, a := range t.Attributes {
		g.out("%s %s", goName(a.Name), columnType(a))
	}
	g.out("}")
	return nil
}

func (g *generator) generateRelations(t *er.EntityType) error {
	return nil
}

func (g *generator) generateIO(t *er.EntityType) error {
	if t.DependsOn != nil {
		g.out("func (t *%s_Set) parse(p *rsf.Reader, parent *%s) error {", goName(t.Name), goName(t.DependsOn.Target.Name))
	} else {
		g.out("func (t *%s_Set) parse(p *rsf.Reader) error {", goName(t.Name))
	}
	g.out("var e %s", goName(t.Name))
	g.out("for p.Next() == rsf.Attribute {")
	g.out("var err error")
	g.out("switch p.Name {")
	for _, a := range t.Attributes {
		g.out("case %q:", a.Name)
		g.out("e.%s, err = %s", goName(a.Name), attrParse(a.Type))
	}
	g.out("default:err= er.UnknownAttribute")
	g.out("}")
	g.out("if err != nil { return err }")
	g.out("}")
	g.out("for p.Next() == rsf.RecordStart {")
	g.out("var err error")
	g.out("switch p.Name {")
	for _, t := range g.dependants(t) {
		g.out("case %q:", t.Name)
		g.out("err = m.%s.parse(p, e)", goName(t.Name))
	}
	g.out("default:err= er.UnknownEntityType")
	g.out("}")
	g.out("if err != nil { return err }")
	g.out("if p.Next() != rsf.RecordEnd {return er.SyntaxError}")
	g.out("}")
	g.out("return nil")
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

func attrType(t er.AttributeType) string {
	switch t {
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
		tn = "er.Int"
	case er.FloatType:
		tn = "er.Float64"
	case er.StringType:
		tn = "er.String"
	default:
		return "?"
	}
	if a.Identifying {
		return tn + "Index"
	}
	return tn + "Column"
}

func attrParse(t er.AttributeType) string {
	switch t {
	case er.IntType:
		return "strconv.Atoi(p.Value)"
	case er.FloatType:
		return "strconv.ParseFloat(p.Value)"
	case er.StringType:
		return "p.Value, nil"
	}
	return "?"
}
