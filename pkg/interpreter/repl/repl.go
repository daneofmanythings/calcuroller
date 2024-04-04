package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/daneofmanythings/calcuroller/pkg/interpreter/evaluator"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/lexer"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/object"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/parser"
)

func run(input string) (object.Object, *object.Metadata) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	value, metadata := evaluator.EvalFromRequest(program)

	return value, metadata
}

func RunFromTerminal() {
	fmt.Println("Welcome to the calcuroller REPL!")
	fmt.Print("(enter dice strings, ex: d20 + 4)\n\n")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		if err == nil {
			val, _ := run(input)
			fmt.Printf("%d\n\n", val.(*object.Integer).Value)
		} else {
			fmt.Printf("\nan error occurred reading input. err=%s", err)
		}
	}
}

func RunFromGRPC(input string) (string, *object.Metadata) {
	value, metadata := run(input)

	valueAsString := fmt.Sprintf("%d", value)

	return valueAsString, metadata
}
