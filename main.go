package main

import (
	"flag"
	"os"
)

var (
	cfgName = flag.String("c", ".builder.yaml", "builder config file")
)

func main() {
	flag.Parse()
	if len(os.Args) == 2 && os.Args[1] == "init" {
		DefaultBuilder()
		return
	}
	eng, err := NewEngine(*cfgName)
	if err != nil {
		printExit(1, err.Error())
	}
	if err = eng.Execute(flag.Args()); err != nil {
		printExit(1, err.Error())
	}
}
