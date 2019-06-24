package gen

import (
	"bytes"
	"github.com/bobappleyard/er"
	"go/format"
	"strings"
	"text/template"
)

func generate(m *er.EntityModel) ([]byte, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"goName":     goName,
		"columnType": columnType,
		"columnInit": columnInit,
		"attrType":   attrType,
		"attrParse":  attrParse,
		"dependants": func(t *er.EntityType) []*er.EntityType {
			return dependants(m, t)
		},
		"newAttrs": newAttrs,
	}).ParseFiles("pkg.tmpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "pkg.tmpl", m); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

func dependants(m *er.EntityModel, t *er.EntityType) []*er.EntityType {
	var res []*er.EntityType
	for _, u := range m.Types {
		var p *er.EntityType
		if u.Dependency.Rel != nil {
			p = u.Dependency.Rel.Target
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

func columnInit(a *er.Attribute) string {
	init := "Column"
	if a.Identifying {
		init = "Index"
	}
	return columnType(a) + init
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

func newAttrs(t *er.EntityType) []*er.Attribute {
	if t.Dependency.Rel == nil {
		return t.Attributes
	}
	var res []*er.Attribute
	for _, a := range t.Attributes {
		if attrIsNew(t, a) {
			res = append(res, a)
		}
	}
	return res
}

func attrIsNew(t *er.EntityType, a *er.Attribute) bool {
	for _, b := range t.Dependency.Rel.Implementation {
		if a == b.Source {
			return false
		}
	}
	return true
}
