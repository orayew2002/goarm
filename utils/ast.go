package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

// AppendFieldStruct adds a new field to a struct in a Go source file.
// `filePath`: path to the .go file
// `structName`: the name of the struct to modify
// `field`: a single field like "Age int" or "Email string `json:\"email\"`"
func AppendFieldStruct(filePath, structName, field string) error {
	// Step 1: Create the file set and parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// Step 2: Split the field string into name and type
	parts := strings.Fields(field)
	if len(parts) < 2 {
		return fmt.Errorf("invalid field format: %q", field)
	}

	fieldName := parts[0]
	fieldType := strings.Join(parts[1:], " ")

	// Step 3: Construct the field using AST
	newField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fieldName)},
		Type:  ast.NewIdent(fieldType),
	}

	// Step 4: Traverse AST and find the struct
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Append the new field
			structType.Fields.List = append(structType.Fields.List, newField)
			found = true
			return false
		}

		return true
	})

	if !found {
		return fmt.Errorf("struct %q not found in %s", structName, filePath)
	}

	// Step 5: Write the modified file back
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if err := printer.Fprint(file, fset, node); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	return nil
}

// AddImportToFile adds a new import path to the import section of a Go file.
// If the import already exists, it does nothing.
func AddImportToFile(filePath string, importPath string) error {
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// Check if import already exists
	for _, imp := range node.Imports {
		if imp.Path.Value == fmt.Sprintf("%q", importPath) {
			// Already imported
			return nil
		}
	}

	// Create a new import spec
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%q", importPath),
		},
	}

	// Ensure there is an import decl block
	importAdded := false
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, newImport)
			importAdded = true
			break
		}
	}

	// If no import block, create one at the top
	if !importAdded {
		newDecl := &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				newImport,
			},
		}
		node.Decls = append([]ast.Decl{newDecl}, node.Decls...)
	}

	// Open file for writing
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer outFile.Close()

	// Write the modified AST back to the file
	if err := printer.Fprint(outFile, fset, node); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
