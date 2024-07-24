package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"script/ast"
	"script/compiler"
	"script/lexer"
	"script/vm"
)

//go:embed example.ys
var example string

const debug = true

func main() {
	//fmt.Println("### Script ###")
	//reader := bufio.NewReader(os.Stdin)

	v := vm.New()

	//for {
	//	fmt.Print("> ")
	//	input, err := reader.ReadString('\n')
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	tokens, errs := lexer.Tokenize([]byte(example))
	if len(errs) > 0 {
		fmt.Printf("error: %+v\n", errs)
		os.Exit(1)
	}
	if debug {
		fmt.Println("tokens:", tokens)
	}

	p, errs := ast.Parse(tokens)
	if len(errs) > 0 {
		fmt.Printf("error: %+v\n", errs)
		os.Exit(1)
	}

	if debug {
		fmt.Println(p)
		fmt.Println("### BYTECODE ###")
	}

	bytecode := make(vm.Bytecode, 0)
	err := compiler.Compile(&bytecode, p)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	if debug {
		fmt.Println(bytecode.String())

		file, err := os.Create("dump.yasm")
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(bytecode.String())
		file.Close()

		fmt.Println("### VM ###")
	}

	err = v.Execute(bytecode)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	if debug {
		fmt.Println(v.Dump())
	}

	//}
}
