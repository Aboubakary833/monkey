package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/internal/lexer"
	"monkey/internal/token"
	"time"
)

const PROMPT = ">> "

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)

	printDatetime()
	fmt.Printf("Type %q for more information.\n", "help")

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)

		for _token := lex.NextToken(); _token.Type != token.EOF; _token = lex.NextToken() {
			fmt.Printf("%v\n", _token)
		}
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
