package main

import (
	"os"

	compiler "mygo/compiler"
)

func main() {
	file, err := os.OpenFile("./sample/sample.mygo", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	parser := compiler.NewParser(file)
	tree := parser.Parse()

	outfile, err := os.OpenFile("./sample/sample.s", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	gen := compiler.NewCgen(tree, outfile)
	gen.GenAST()




}
