package main

import (
	"monkey/cmd/repl"
	"os"
)

func main() {
	
	repl.Start(os.Stdin, os.Stdout)

}
