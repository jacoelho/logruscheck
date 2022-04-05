package logruscheck

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"regexp"
)

var snakeCaseRE = regexp.MustCompile("^[a-z]+(_[a-z]+)*$")

var Analyzer = &analysis.Analyzer{
	Name:     "logruscheck",
	Doc:      "Checks that log keys are in snake_case and structured logging is used",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		fun, ok := funcForCallExpr(pass, call)
		if !ok {
			return
		}

		if !isFromPkg(fun, "github.com/sirupsen/logrus") {
			return
		}

		switch fun.Name() {
		case "WithField":
			isWithFieldFnCall(pass, call)
		case "WithFields":
			isWithFieldsFnCall(pass, call)
		case "Tracef", "Debugf", "Printf", "Infof", "Warnf", "Warningf", "Errorf", "Panicf", "Fatalf":
			isLogFormatFnCall(pass, call)
		}
	})

	return nil, nil
}

func isLogFormatFnCall(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	funcDecl, ok := pass.TypesInfo.TypeOf(sel.Sel).(*types.Signature)
	if !ok {
		return
	}

	// [0] must be format (string)
	// [1] must be args (...interface{})
	if funcDecl.Params().Len() != 2 {
		return
	}

	firstArg, ok := funcDecl.Params().At(0).Type().(*types.Basic)
	if !ok {
		return
	}

	if firstArg.Kind() != types.String {
		return
	}

	secondArg, ok := funcDecl.Params().At(1).Type().(*types.Slice)
	if !ok {
		return
	}

	if _, ok := secondArg.Elem().(*types.Interface); !ok {
		return
	}

	reportLoggingFormat(pass, sel.Sel, sel.Sel.Name)
}

// isWithFieldFnCall log.WithField("foo", "bar")
func isWithFieldFnCall(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}

	firstArg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	if firstArg.Kind != token.STRING {
		return
	}

	trimmed := firstArg.Value[1 : len(firstArg.Value)-1]
	if trimmed == "" {
		return
	}

	if snakeCaseRE.MatchString(trimmed) {
		return
	}

	reportSnakeCase(pass, firstArg, trimmed)
}

// isWithFieldsFnCall
func isWithFieldsFnCall(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) != 1 {
		return
	}

	firstArg, ok := call.Args[0].(*ast.CompositeLit)
	if !ok {
		return
	}

	m, ok := pass.TypesInfo.TypeOf(firstArg.Type).Underlying().(*types.Map)
	if !ok {
		return
	}

	if m.String() != "map[string]interface{}" && m.String() != "map[string]any" {
		return
	}

	for _, expr := range firstArg.Elts {
		kv, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		k, ok := kv.Key.(*ast.BasicLit)
		if !ok {
			continue
		}

		if k.Kind != token.STRING {
			continue
		}

		trimmed := k.Value[1 : len(k.Value)-1]
		if trimmed == "" {
			continue
		}

		if snakeCaseRE.MatchString(trimmed) {
			continue
		}

		reportSnakeCase(pass, k, trimmed)
	}
}

func funcForCallExpr(pass *analysis.Pass, call *ast.CallExpr) (*types.Func, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	// Get the object this is calling.
	obj, ok := pass.TypesInfo.Uses[sel.Sel]
	if !ok {
		return nil, false
	}

	f, ok := obj.(*types.Func)
	return f, ok
}

func isFromPkg(fun *types.Func, pkg string) bool {
	if fun.Pkg() == nil {
		return false
	}

	return fun.Pkg().Path() == pkg
}

func reportSnakeCase(pass *analysis.Pass, expr ast.Expr, word string) {
	pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		End:            expr.End(),
		Category:       "logging",
		Message:        fmt.Sprintf("log key '%s' should be snake_case.", word),
		SuggestedFixes: nil,
	})
}

func reportLoggingFormat(pass *analysis.Pass, expr ast.Expr, word string) {
	pass.Report(analysis.Diagnostic{
		Pos:            expr.Pos(),
		End:            expr.End(),
		Category:       "logging",
		Message:        fmt.Sprintf("call to '%s' should be replaced with WithField or WithFields.", word),
		SuggestedFixes: nil,
	})
}
