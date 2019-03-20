package gen

import (
	"github.com/bobappleyard/er"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

func generate(m *er.EntityModel) (*ast.File, error) {
	f := &ast.File{
		Name: &ast.Ident{Name: m.Name},
	}
	addImports(f, "sort", "github.com/bobappleyard/er")
	fs := addstructType(f, "Model")
	generateNewFunc(f, m)
	generateValidateMethod(f, fs, m)
	for _, t := range m.Types {
		generateType(f, t)
		generateListType(f, fs, t)
	}
	return f, nil
}

func generateNewFunc(f *ast.File, m *er.EntityModel) {
	nimp := addFunc(f, nil, "New", nil, ret(ptr(id("Model"))))
	nimp.List = append(nimp.List, &ast.AssignStmt{
		Lhs: []ast.Expr{id("m")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{call(id("new"), id("Model"))},
	})
	for _, t := range m.Types {
		tn := snakeToGo(t.Name)
		nimp.List = append(nimp.List, assign(sel(id("m"), tn), &ast.CompositeLit{
			Type: id(tn + "_List"),
			Elts: []ast.Expr{&ast.KeyValueExpr{
				Key:   id("model"),
				Value: id("m"),
			}},
		}))
	}
	nimp.List = append(nimp.List, &ast.ReturnStmt{Results: []ast.Expr{
		id("m"),
	}})
}

func generateValidateMethod(f *ast.File, fs structType, m *er.EntityModel) {
	vimp := fs.addPointerMethod("Validate", nil, args(field("err", id("error"))))
	for _, t := range m.Types {
		tn := snakeToGo(t.Name)
		vimp.List = append(vimp.List, &ast.ExprStmt{X: call(sel(id("m"), tn, "validate"), &ast.UnaryExpr{
			Op: token.AND,
			X:  id("err"),
		})})
	}
	vimp.List = append(vimp.List, &ast.ReturnStmt{Results: []ast.Expr{
		id("err"),
	}})
}

func generateType(f *ast.File, t *er.EntityType) {
	tn := snakeToGo(t.Name)

	var attrs []ast.Spec

	for i, a := range t.Attributes {
		attr := &ast.ValueSpec{
			Names: []*ast.Ident{{Name: "Attr_" + tn + "_" + snakeToGo(a.Name)}},
		}
		if i == 0 {
			attr.Values = []ast.Expr{id("iota")}
		}
		attrs = append(attrs, attr)
	}
	f.Decls = append(f.Decls, &ast.GenDecl{
		Lparen: 1,
		Tok:    token.CONST,
		Specs:  attrs,
	})

	tfs := addstructType(f, tn)
	tfs.addField("model", ptr(id("Model")))

	for _, a := range t.Attributes {
		tfs.addField(snakeToGo(a.Name), attributeType(a))
	}
	for _, r := range t.Relationships {
		m := tfs.addMethod(snakeToGo(r.Name), nil, ret(id(snakeToGo(r.Target.Name))))
		generateRelationshipImplementation(f, r, m)
	}
}

func generateRelationshipImplementation(f *ast.File, r *er.Relationship, impb *ast.BlockStmt) {
	targ := snakeToGo(r.Target.Name)
	rn := receiverName(r.Source.Name)
	imp := &ast.CompositeLit{
		Type: sel(id("er"), "Query"),
	}
	s := &ast.ReturnStmt{
		Results: []ast.Expr{
			index(call(
				sel(id(rn), "model", targ, "Lookup"),
				imp,
			), &ast.BasicLit{Kind: token.INT, Value: "0"}),
		},
	}
	impb.List = append(impb.List, s)
	for _, c := range r.Implementation {
		var path ast.Expr = id(rn)
		for _, step := range c.BasePath {
			path = call(sel(path, snakeToGo(step.Rel.Name)))
		}
		path = sel(path, snakeToGo(c.Source.Name))
		imp.Elts = append(imp.Elts, path)
	}
}

func generateListType(f *ast.File, fs structType, t *er.EntityType) {
	tn := snakeToGo(t.Name)
	ltn := tn + "_List"
	fs.addField(tn, id(ltn))
	tfs := addstructType(f, ltn)
	tfs.addField("model", ptr(id("Model")))
	tfs.addField("meta", sel(id("er"), "EntityMeta"))
	tfs.addField("items", &ast.ArrayType{Elt: id(tn)})
	generateLookupMethod(tfs, t)
}

func generateLookupMethod(tfs structType, t *er.EntityType) {
	tn := snakeToGo(t.Name)
	rn := receiverName(tn)
	imp := tfs.addMethod("Lookup", args(field("q", sel(id("er"), "Query"))), args(field("res", &ast.ArrayType{Elt: id(tn)}), field("err", id("error"))))
	imp.List = append(imp.List,
		define(id("idxs"), call(sel(id(rn), "meta", "EvalQuery"), id("q"), sel(id(rn), "items"))),
		&ast.ForStmt{
			Cond: call(sel(id("idxs"), "Next")),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					assign(id("res"), call(id("append"),
						id("res"),
						index(sel(id(rn), "items"), call(sel(id("idxs"), "This"))),
					)),
				},
			},
		},
		&ast.ReturnStmt{Results: []ast.Expr{
			id("res"),
			call(sel(id("idxs"), "Err")),
		}},
	)
}

func attributeType(a *er.Attribute) ast.Expr {
	var at string
	switch a.Type {
	case er.StringType:
		at = "string"
	case er.IntType:
		at = "int"
	case er.FloatType:
		at = "float64"
	}
	return id(at)
}

func attributeColumnType(a *er.Attribute) ast.Expr {
	var at string
	switch a.Type {
	case er.StringType:
		at = "String"
	case er.IntType:
		at = "Int"
	case er.FloatType:
		at = "Float"
	}
	if a.Identifying {
		at += "Index"
	} else {
		at += "Column"
	}
	return sel(id("er"), at)
}

func receiverName(n string) string {
	return strings.ToLower(n[:1])
}

func generateKeyClause(kimp *ast.BlockStmt, kk ast.Expr, a *er.Attribute) {
	name := snakeToGo(a.Name)
	base := receiverName(a.Owner.Name)
	kimp.List = append(kimp.List, &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  sel(id(base), name),
			Op: token.NEQ,
			Y:  sel(kk, name),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.BinaryExpr{
						X:  sel(id(base), name),
						Op: token.LSS,
						Y:  sel(kk, name),
					},
				}},
			},
		},
	})
}

func generateKeyGet(kget *ast.CompositeLit, a *er.Attribute) {
	aname := snakeToGo(a.Name)
	kget.Elts = append(kget.Elts, &ast.KeyValueExpr{
		Key: id(aname),
		Value: &ast.SelectorExpr{
			X:   id(receiverName(a.Owner.Name)),
			Sel: id(aname),
		},
	})
}

var snPat = regexp.MustCompile(`(^|_).`)

func snakeToGo(id string) string {
	return snPat.ReplaceAllStringFunc(id, func(s string) string {
		if s[0] == '_' {
			s = s[1:]
		}
		return strings.ToUpper(s)
	})
}
