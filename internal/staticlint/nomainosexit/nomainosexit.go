// Package nomainosexit содержит анализатор кода, который запрещает использование прямого вызова os.Exit()
// внутри функции main() в пакете main.
package nomainosexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	checkMessage = "direct call of os.Exit() in main() function of main package is not allowed"
)

// Analyzer анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.
var Analyzer = &analysis.Analyzer{
	Name: "nomainosexit",
	Doc:  checkMessage,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// for each passed file
	for _, file := range pass.Files {

		// check only for main package
		if pass.Pkg.Name() != "main" {
			continue
		}

		inMainFunc := false
		ast.Inspect(file, func(node ast.Node) bool {
			// found function declaration?
			if fn, nodeIsFuncDecl := node.(*ast.FuncDecl); nodeIsFuncDecl {

				// function declaration is main?
				inMainFunc = false
				if fn.Name.String() == "main" {
					inMainFunc = true
				}

				return true
			}

			// found function call knowing, that we entered main function?
			if call, nodeIsCallExpr := node.(*ast.CallExpr); nodeIsCallExpr {

				// observed os.Exit() call - reporting
				if inMainFunc && isPackageMethodCall(call, "os", "Exit") {
					pass.Reportf(call.Pos(), checkMessage)
				}

				return true
			}

			return true
		})
	}

	return nil, nil
}

// isPackageMethodCall проверяет, является ли выражение вызовом метода определенного пакета с указанным именем.
func isPackageMethodCall(callExpr *ast.CallExpr, pkg, name string) bool {
	funcExpr := callExpr.Fun

	// expression followed by a selector as in pkg.MethodCall() expression
	selectorExpr, ok := funcExpr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent := selectorExpr.X
	methodIdent := selectorExpr.Sel

	return isIdent(pkgIdent, pkg) && // check pkg name
		isIdent(methodIdent, name) // check
}

// isIdent проверяет, является ли выражение идентификатором с заданным именем.
func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
