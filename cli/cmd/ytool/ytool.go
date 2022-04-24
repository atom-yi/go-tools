package main

import (
	"os"

	"github.com/atom-yi/cli/pkg"
)

func main() {
	pkg.PrintEvn(os.Args[1])
}
