package protobuf

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func addGoMethodToService(filePath, method string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	funcToAdd := fmt.Sprintf(`
func (s *Service) %s(ctx context.Context, c *connect.Request[gen.%sRequest]) (*connect.Response[gen.%sResponse], error) {
    return connect.NewResponse(&gen.%sResponse{}), nil
}
`, method, method, method, method)

	funcNode, err := parser.ParseFile(fset, "", "package p\n"+funcToAdd, 0)
	if err != nil {
		return fmt.Errorf("error parsing function to add: %w", err)
	}
	funcDecl := funcNode.Decls[0].(*ast.FuncDecl) // Extract the function declaration

	found := false
	for _, d := range node.Decls {
		if genDecl, ok := d.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok && typeSpec.Name.Name == "Service" {
						// Add the function after the "type Service struct" declaration
						node.Decls = append(node.Decls, funcDecl)
						found = true
						break
					}
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("no 'type Service struct' declaration found in the file")
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("error formatting modified AST: %w", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("error writing back to file: %w", err)
	}
	return nil
}
