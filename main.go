package main

import (
	"gamefed/srv/log"
)

func main() {
	shout()

	log := logsrv.Begin()
	defer log.End()

	log.Highlightf("gamefed ⚡")
	defer log.Highlightf("gamefed 👻")

	run()
}
