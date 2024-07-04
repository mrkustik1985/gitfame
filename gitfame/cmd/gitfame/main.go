//go:build !solution

package main

import (
	"github.com/sirupsen/logrus"
	parser2 "gitlab.com/slon/shad-go/gitfame/pkg/parser"
	"gitlab.com/slon/shad-go/gitfame/pkg/scaner"
	"os"
)

var Scaner scaner.Scaner
var Log *logrus.Logger

func main() {
	Log = logrus.New()
	Log.SetLevel(logrus.DebugLevel)
	//Log.Debug("start parse")
	args := os.Args[1:]
	Scaner.Scan(args)
	//Log.Debug("start routine")
	parser := parser2.NewParser(&Scaner)
	errs := parser.DoRoutine()
	//Log.Debug("finish routine")
	if errs != nil {
		Log.Fatal(errs)
	}
	formatter, err := parser2.NewFormatter(parser.Scaner.Format, parser.Scaner.OrderBy)
	if err != nil {
		Log.Fatal(err)
	}
	err = formatter.Output(parser.Stats)
	if err != nil {
		Log.Fatal(err)
	}
}
