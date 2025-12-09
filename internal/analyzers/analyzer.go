package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "staticlint",
	Doc:  "Static lint checks that static analysis failures",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		callExpr := node.(*ast.CallExpr)

		fun, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		ident, ok := fun.X.(*ast.Ident)
		if !ok {
			return
		}

		if ident.Name == "os" && fun.Sel.Name == "Exit" {
			for _, file := range pass.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					if fn, ok := n.(*ast.FuncDecl); ok {
						if fn.Name.Name == "main" && file.Name.Name == "main" {
							if callExpr.Pos() >= fn.Pos() && callExpr.Pos() <= fn.End() {
								pass.Reportf(
									callExpr.Pos(), "os.Exit called in main function",
								)
							}
						}
					}
					return true
				})
			}
		}
	})
	return nil, nil
}
