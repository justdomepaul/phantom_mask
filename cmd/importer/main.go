package main

import (
	"github.com/justdomepaul/toolbox/errorhandler"
)

var (
	system = "Import Data"
)

func main() {
	defer errorhandler.PanicErrorHandler(system, "import data interrupt => \n")

	_, cleanup, err := Runner()
	if err != nil {
		panic(err)
	}
	defer cleanup()
}
