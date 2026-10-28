package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: check <registry root>")
		return 2
	}
	problems := Check(os.DirFS(args[0]))
	if len(problems) == 0 {
		return 0
	}
	// File names and keys come from the pull request, so the runner must not read them as workflow commands.
	token := rand.Text()
	fmt.Fprintf(stdout, "::stop-commands::%s\n", token)
	for _, p := range problems {
		fmt.Fprintln(stdout, p)
	}
	fmt.Fprintf(stdout, "::%s::\n", token)
	return 1
}
