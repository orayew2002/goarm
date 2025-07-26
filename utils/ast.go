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

// AppendFuncArgument adds a new argument to the specified function in the Go source file.
func AppendFuncArgument(filePath, funcName, newArgName, newArgType string) error {
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	found := false

	// Walk through the AST
	for _, decl := range node.Decls {
		// We only want function declarations
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != funcName {
			continue
		}

		// Create a new field (argument)
		newField := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(newArgName)},
			Type:  parseExprFromType(newArgType),
		}

		// Append the argument
		funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, newField)
		found = true
		break
	}

	if !found {
		return fmt.Errorf("function %q not found", funcName)
	}

	// Open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	// Write modified AST to file
	if err := printer.Fprint(file, fset, node); err != nil {
		return fmt.Errorf("failed to write modified file: %w", err)
	}

	return nil
}

// parseExprFromType parses a Go type (e.g., "context.Context") into an AST expression.
func parseExprFromType(typ string) ast.Expr {
	// Support for qualified types like context.Context
	parts := strings.Split(typ, ".")
	if len(parts) == 2 {
		return &ast.SelectorExpr{
			X:   ast.NewIdent(parts[0]),
			Sel: ast.NewIdent(parts[1]),
		}
	}
	// Simple types like int, string
	return ast.NewIdent(typ)
}

// AddReturnFieldToConstructor modifies the return expression of a constructor
// (e.g. NewRepo) to include a struct literal field initialization (e.g. db: db).
func AddReturnFieldToConstructor(filePath, funcName, fieldName string) error {
	fset := token.NewFileSet()

	// Parse Go file
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	// Traverse the AST and find the target function
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != funcName {
			return true
		}

		// Find the return statement inside the function body
		for _, stmt := range fn.Body.List {
			ret, ok := stmt.(*ast.ReturnStmt)
			if !ok || len(ret.Results) == 0 {
				continue
			}

			compLit, ok := ret.Results[0].(*ast.UnaryExpr)
			if !ok || compLit.Op != token.AND {
				continue
			}

			structLit, ok := compLit.X.(*ast.CompositeLit)
			if !ok {
				continue
			}

			// Add field to composite literal
			structLit.Elts = append(structLit.Elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(fieldName),
				Value: ast.NewIdent(fieldName),
			})

			break
		}
		return false
	})

	// Write the modified file back
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return printer.Fprint(file, fset, node)
}

// AddArgumentToFunctionCall adds an argument to the specified function call (e.g., repo.NewRepo)
func AddArgumentToFunctionCall(filePath, fullFuncName, argName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	// Split the function name: "repo.NewRepo" => "repo", "NewRepo"
	parts := strings.Split(fullFuncName, ".")
	if len(parts) != 2 {
		return nil // Invalid format
	}
	pkgName, funcName := parts[0], parts[1]

	// Search for and modify the desired function call
	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if this is a call to repo.NewRepo
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
				if pkgIdent.Name == pkgName && selExpr.Sel.Name == funcName {
					// Add argument only if the call currently has no arguments
					if len(callExpr.Args) == 0 {
						callExpr.Args = []ast.Expr{
							ast.NewIdent(argName),
						}
					}
				}
			}
		}
		return true
	})

	// Overwrite the file with modified AST
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return printer.Fprint(outFile, fset, node)
}
