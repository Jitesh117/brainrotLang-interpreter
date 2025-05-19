package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Jitesh117/brainrotLang-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("No cap %s! This is the BrainrotLang programming language!\n",
		user.Username)
	fmt.Printf("It's giving ✨runtime✨")
	fmt.Printf("(hit that Ctrl+C once it gets cringe though)\n")

	repl.Start(os.Stdin, os.Stdout)
}
