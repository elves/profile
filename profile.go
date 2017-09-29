package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/elves/elvish/daemon/api"
	"github.com/elves/elvish/eval"
	"github.com/elves/elvish/parse"
)

var sources = map[string]string{
	"0":   "range 5000000 | each [_]{ }",
	"1":   "range 5000000 | each [_]{ nop }",
	"2":   "range 5000000 | each $&nop",
	"for": "for x [(range 1000)] { for y [(range 1000)] { } }",
}

func assertOK(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func profile(name string) {
	sourceBytes, err := ioutil.ReadFile(name)
	assertOK(err)
	source := string(sourceBytes)
	profileFile, err := os.Create(name + ".prof")
	assertOK(err)
	defer profileFile.Close()
	ast, err := parse.Parse(name, source)
	assertOK(err)
	ev := eval.NewEvaler(api.NewClient("/invalid"), nil, "", nil)
	op, err := ev.Compile(ast, name, source)

	pprof.StartCPUProfile(profileFile)
	defer pprof.StopCPUProfile()
	ev.Eval(op, name, source)
}

func main() {
	for _, name := range os.Args[1:] {
		profile(name)
	}
}
