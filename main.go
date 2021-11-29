package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"terminascript/evaluator"
	"terminascript/lexer"
	"terminascript/parser"
)

func ReadFile(filename string) string {
	filePointer, _ := os.Open(filename)
	fileBytes, _ := ioutil.ReadAll(filePointer)
	return string(fileBytes)
}

func startRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, ">>")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		interpretProgram(line)
	}
}

func interpretProgram(program string) {
	l := lexer.NewLexer(strings.TrimSpace(program))
	tokens := l.Lex()

	p := parser.NewParser(tokens)
	ast := p.Parse()

	fmt.Println(ast)
	evaluator.Eval(ast)
}

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		file := ReadFile(filename)
		formattedFile := strings.Replace(file, `\n`, ``, -1)
		interpretProgram(formattedFile)
	} else {
		startRepl(os.Stdin, os.Stdout)
	}
}
