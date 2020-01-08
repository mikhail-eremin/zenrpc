package parser_v2

import (
	"fmt"
	"go/ast"
)

type Service struct {
	Name        string
	Description string

	Methods []*Method
}

func NewService(spec *ast.TypeSpec) *Service {
	return &Service{
		Name:        spec.Name.Name,
		Description: parseComment(spec.Doc),

		Methods: []*Method{},
	}
}

func (s *Service) AddMethod(method *Method) {
	s.Methods = append(s.Methods, method)
}

func hasMagicComment(spec *ast.TypeSpec) bool {
	if spec.Comment != nil && len(spec.Comment.List) > 0 && spec.Comment.List[0].Text == MagicComment {
		return true
	}

	return false
}

func hasEmbedService(spec *ast.TypeSpec) bool {
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return false
	}

	if structType.Fields.List == nil {
		return false
	}

	for _, field := range structType.Fields.List {
		expr, ok := field.Type.(*ast.SelectorExpr)
		if !ok || expr.Sel == nil {
			continue
		}

		x, ok := expr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if fmt.Sprintf("%s.%s", x.Name, expr.Sel.Name) == ServiceType {
			return true
		}
	}

	return false
}
