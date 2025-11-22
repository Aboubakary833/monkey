package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/internal/lexer"
	"monkey/internal/parser"
	"time"
)

const PROMPT = ">> "
const HELP_COMMAND = "help"

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	printDatetime()
	fmt.Printf("Type %q for more information.\n", HELP_COMMAND)

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		parser := parser.New(lex)
		
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			for _, errMsg := range parser.Errors() {
				io.WriteString(output, errMsg)
				io.WriteString(output, "\n")
			}
			continue
		}

		io.WriteString(output, program.String())
		io.WriteString(output, "\n")
	}
}

func printDatetime() {
	localTime := time.Now().Local()
	year, month, day := localTime.Date()
	
	hour := localTime.Hour()
	min := localTime.Minute()
	sec := localTime.Second()

	fmt.Printf(
		"Welcome to Monkey 1.0.0 beta (main - %v %d %d, %d:%d:%d)\n",
		month, day, year, hour, min, sec,
	)
}
