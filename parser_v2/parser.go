package parser_v2

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
)

const (
	MagicComment = `//zenprc`
	ServiceType  = `zenrpc.Service`
)

type Parser struct {
	entryPoint string

	pkgs []*packages.Package

	services map[string]*Service
}

func NewParser(filename string) *Parser {
	return &Parser{
		entryPoint: filename,

		services: map[string]*Service{},
	}
}

func (p *Parser) Parse() error {
	if err := p.loadPackages(); err != nil {
		return err
	}

	p.loadServices()
	p.loadMethods()

	return nil
}

func (p *Parser) Service(name string) *Service {
	return p.services[name]
}

func (p *Parser) loadPackages() error {
	fi, err := os.Stat(p.entryPoint)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%s is directory", p.entryPoint)
	}

	config := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	p.pkgs, err = packages.Load(config, p.entryPoint)

	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) loadServices() {
	for _, pkg := range p.pkgs {
		for _, stx := range pkg.Syntax {
			for _, decl := range stx.Decls {
				gdecl, ok := decl.(*ast.GenDecl)

				if !ok || gdecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range gdecl.Specs {
					spec, ok := spec.(*ast.TypeSpec)
					if !ok || !ast.IsExported(spec.Name.Name) {
						continue
					}

					// check that struct is our zenrpc struct
					if hasMagicComment(spec) || hasEmbedService(spec) {
						service := NewService(spec)
						p.services[service.Name] = service
					}
				}
			}
		}
	}
}

func (p *Parser) loadMethods() {
	for _, pkg := range p.pkgs {
		for _, stx := range pkg.Syntax {
			for _, decl := range stx.Decls {
				fdecl, ok := decl.(*ast.FuncDecl)
				if !ok || fdecl.Recv == nil {
					continue
				}

				method := NewMethod(fdecl)
				service := p.Service(method.ServiceName)

				if service == nil {
					continue
				}

				service.AddMethod(method)
			}
		}
	}
}

func parseComment(comment *ast.CommentGroup) string {
	if comment == nil {
		return ""
	}

	result := ""
	for _, comment := range comment.List {
		if strings.HasPrefix(comment.Text, MagicComment) {
			continue
		}

		if len(result) > 0 {
			result += "\n"
		}

		result += strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
	}

	return result
}
