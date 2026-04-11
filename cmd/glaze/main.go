package main

import (
	"fmt"
	"os"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/evaluator"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/lexer"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/object"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/parser"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input file provided.")
		fmt.Println("Usage: ./glaze <filename>")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Multiple files found.")
		fmt.Println("Usage: ./glaze <filename>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	source := string(data)

	// LEXING (Text -> Tokens)
	s := lexer.InitScanner(source)

	var tokens []token.Token
	for {
		tok := s.ScanToken()
		tokens = append(tokens, tok)
		if tok.Type == token.END_OF_FILE {
			break
		}
	}

	// PARSING (Tokens -> AST)
	p := parser.New(tokens)
	program := p.ParseProgram()

	// Stop execution if there are syntax errors!
	if len(p.Errors()) != 0 {
		fmt.Println("PARSER ERRORS:")
		for _, msg := range p.Errors() {
			fmt.Printf("\t- %s\n", msg)
		}
		os.Exit(1)
	}

	// EVALUATING (AST -> Execution)
	env := object.NewEnvironment()
	env.Define("bool", &object.String{Value: "bool"}, true, "STRING", true)
	env.Define("char", &object.String{Value: "char"}, true, "STRING", true)
	env.Define("i8", &object.String{Value: "i8"}, true, "STRING", true)
	env.Define("i16", &object.String{Value: "i16"}, true, "STRING", true)
	env.Define("i32", &object.String{Value: "i32"}, true, "STRING", true)
	env.Define("i64", &object.String{Value: "i64"}, true, "STRING", true)

	env.Define("u8", &object.String{Value: "u8"}, true, "STRING", true)
	env.Define("u16", &object.String{Value: "u16"}, true, "STRING", true)
	env.Define("u32", &object.String{Value: "u32"}, true, "STRING", true)
	env.Define("u64", &object.String{Value: "u64"}, true, "STRING", true)

	env.Define("f32", &object.String{Value: "f32"}, true, "STRING", true)
	env.Define("f64", &object.String{Value: "f64"}, true, "STRING", true)
	result := evaluator.Eval(program, env)

	// Catch and print any runtime type-mismatches or undeclared variable errors
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println("RUNTIME ERROR:")
		fmt.Println(result.Inspect())
		os.Exit(1)
	}
}
