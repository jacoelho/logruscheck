package main

import (
	"github.com/jacoelho/logruscheck/logruscheck"
	"golang.org/x/tools/go/analysis"
)

// analyzerPlugin implements the golangci-lint AnalyzerPlugin interface.
// see https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
type analyzerPlugin struct{}

// AnalyzerPlugin is the golangci-lint plugin.
var AnalyzerPlugin analyzerPlugin //nolint:deadcode,unused

// GetAnalyzers returns all analyzers for a plugin
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		logruscheck.Analyzer,
	}
}
