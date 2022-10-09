package main

import (
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/shutdown"
	"os"
)

var (
	system = "Restful Server"
)

func main() {
	defer errorhandler.PanicErrorHandler(system, "Restful server interrupt => \n")

	_, cleanup, err := RestfulRunner()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	quit := make(chan os.Signal)
	defer close(quit)
	shutdown.NewShutdown(
		shutdown.WithQuit(quit),
	).Shutdown()
}
