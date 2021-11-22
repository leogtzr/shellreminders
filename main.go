package main

import (
	"os"

	"github.com/leogtzr/shellreminders/shellreminders"
)

func main() {
	os.Exit(shellreminders.CLI(os.Args[:]))
}
