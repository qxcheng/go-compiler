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

	scanner := compiler.NewScanner(file)
	token, _ := scanner.GetToken()
	for token != compiler.ENDFILE {
		token, _ = scanner.GetToken()
	}

}
