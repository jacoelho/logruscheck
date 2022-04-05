package main

import (
	"github.com/jacoelho/logruscheck/logruscheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logruscheck.Analyzer)
}
