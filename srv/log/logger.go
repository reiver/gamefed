package logsrv

import (
	"io"
	"os"

	"codeberg.org/reiver/go-log"

	"gamefed/cfg"
)

var writer io.Writer = os.Stdout

var logger log.Logger = log.CreateLogger(writer, cfg.LogLevel())
