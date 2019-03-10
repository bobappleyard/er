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
	fs := addStructType(f, "Model")
	for _, t := range m.Types {
		generateType(f, t)
		generateListType(f, fs, t)
	}
	return f, nil
}

func generateType(f *ast.File, t *er.EntityType) {
	tn := snakeToGo(t.Name)
	tfs := addStructType(f, tn)
	addField(tfs, "model", &ast.StarExpr{X: id("Model")})

	kfs := addStructType(f, tn+"_Key")
	kimp := addMethod(f, field("k", id(tn+"_Key")), "keyCompare",
		[]*ast.Field{field("to", id(tn+"_Key"))},
		[]*ast.Field{{Type: id("bool")}})

	for _, a := range t.Attributes {
		addField(tfs, snakeToGo(a.Name), attributeType(a))
		if a.Identifying {
			addField(kfs, snakeToGo(a.Name), attributeType(a))
			generateKeyClause(kimp, a)
		}
	}
	kimp.List = append(kimp.List, &ast.ReturnStmt{
		Results: []ast.Expr{id("true")},
	})
	for _, r := range t.Relationships {
		addMethod(f, field("e", id(tn)), snakeToGo(r.Name), nil, []*ast.Field{{Type: id(snakeToGo(r.Target.Name))}})
	}
}

func generateListType(f *ast.File, fs *ast.FieldList, t *er.EntityType) {
	tn := snakeToGo(t.Name)
	ltn := tn + "_List"
	addField(fs, tn, id(ltn))
	tfs := addStructType(f, ltn)
	addField(tfs, "items", &ast.ArrayType{Elt: id(tn)})
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

func generateKeyClause(kimp *ast.BlockStmt, a *er.Attribute) {
	name := snakeToGo(a.Name)
	kimp.List = append(kimp.List, &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.SelectorExpr{X: id("k"), Sel: id(name)},
			Op: token.NEQ,
			Y:  &ast.SelectorExpr{X: id("to"), Sel: id(name)},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{
					&ast.BinaryExpr{
						X:  &ast.SelectorExpr{X: id("k"), Sel: id(name)},
						Op: token.LSS,
						Y:  &ast.SelectorExpr{X: id("to"), Sel: id(name)},
					},
				}},
			},
		},
	})
}

func id(n string) *ast.Ident {
	return &ast.Ident{Name: n}
}

func addField(fs *ast.FieldList, n string, t ast.Expr) {
	fs.List = append(fs.List, field(n, t))
}

func field(n string, t ast.Expr) *ast.Field {
	return &ast.Field{Names: []*ast.Ident{id(n)}, Type: t}
}

func addMethod(f *ast.File, recv *ast.Field, n string, args, res []*ast.Field) *ast.BlockStmt {
	body := &ast.BlockStmt{}
	f.Decls = append(f.Decls, &ast.FuncDecl{
		Recv: &ast.FieldList{List: []*ast.Field{recv}},
		Name: id(n),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: args},
			Results: &ast.FieldList{List: res},
		},
		Body: body,
	})
	return body
}

func addType(f *ast.File, name string, t ast.Expr) {
	f.Decls = append(f.Decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: id(name),
				Type: t,
			},
		},
	})
}

func addStructType(f *ast.File, name string) *ast.FieldList {
	fields := &ast.FieldList{}
	addType(f, name, &ast.StructType{Fields: fields})
	return fields
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
