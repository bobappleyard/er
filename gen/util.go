package gen

import (
	"go/ast"
	"go/token"
	"strconv"
)

func id(n string) *ast.Ident {
	return &ast.Ident{Name: n}
}

func ptr(t ast.Expr) ast.Expr {
	return &ast.StarExpr{X: t}
}

func field(n string, t ast.Expr) *ast.Field {
	return &ast.Field{Names: []*ast.Ident{id(n)}, Type: t}
}

func args(fs ...*ast.Field) []*ast.Field {
	return fs
}

func ret(t ast.Expr) []*ast.Field {
	return []*ast.Field{{Type: t}}
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

type structType struct {
	name string
	f    *ast.File
	fs   *ast.FieldList
}

func addstructType(f *ast.File, name string) structType {
	fields := &ast.FieldList{}
	addType(f, name, &ast.StructType{Fields: fields})
	return structType{name, f, fields}
}

func (t structType) addField(n string, ft ast.Expr) {
	t.fs.List = append(t.fs.List, field(n, ft))
}

func addFunc(f *ast.File, recv *ast.Field, name string, in, out []*ast.Field) *ast.BlockStmt {
	body := &ast.BlockStmt{}
	var r *ast.FieldList
	if recv != nil {
		r = &ast.FieldList{List: []*ast.Field{recv}}
	}
	d := &ast.FuncDecl{
		Recv: r,
		Name: id(name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: in},
			Results: &ast.FieldList{List: out},
		},
		Body: body,
	}
	f.Decls = append(f.Decls, d)
	return body
}

func (t structType) addMethod(name string, in, out []*ast.Field) *ast.BlockStmt {
	recv := field(receiverName(t.name), id(t.name))
	return addFunc(t.f, recv, name, in, out)
}

func (t structType) addPointerMethod(name string, in, out []*ast.Field) *ast.BlockStmt {
	recv := field(receiverName(t.name), ptr(id(t.name)))
	return addFunc(t.f, recv, name, in, out)
}

func addImports(f *ast.File, names ...string) {
	var imps []ast.Spec
	for _, n := range names {
		imps = append(imps, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(n),
			},
		})
	}
	f.Decls = append(f.Decls, &ast.GenDecl{
		Tok:    token.IMPORT,
		Specs:  imps,
		Lparen: 1,
	})
}

func call(f ast.Expr, args ...ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  f,
		Args: args,
	}
}

func sel(x ast.Expr, s ...string) ast.Expr {
	for _, s := range s {
		x = &ast.SelectorExpr{
			X:   x,
			Sel: id(s),
		}
	}
	return x
}

func index(x, idx ast.Expr) ast.Expr {
	return &ast.IndexExpr{
		X:     x,
		Index: idx,
	}
}

func assign(lhs ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}

func define(lhs ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}
