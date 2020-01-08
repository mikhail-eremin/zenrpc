package parser_v2

import "go/ast"

type Method struct {
	FuncDecl      *ast.FuncType
	Name          string
	ServiceName   string
	LowerCaseName string
	HasContext    bool
	//Args          []Arg
	//DefaultValues map[string]DefaultValue
	//Returns       []Return
	//SMDReturn     *SMDReturn // return for generate smd schema; pointer for nil check
	Description string

	//Errors []SMDError // errors for documentation in SMD
}

func NewMethod(decl *ast.FuncDecl) *Method {
	return &Method{
		Name:        decl.Name.Name,
		ServiceName: parseReceiver(decl),
		Description: parseComment(decl.Doc),
	}
}

func parseReceiver(decl *ast.FuncDecl) string {
	for _, field := range decl.Recv.List {
		// field can be pointer or not
		switch v := field.Type.(type) {
		case *ast.StarExpr:
			if ident, ok := v.X.(*ast.Ident); ok {
				return ident.Name
			}
			continue
		case *ast.Ident:
			return v.Name
		}
	}

	return ""
}
