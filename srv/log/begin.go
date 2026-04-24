package logsrv

import (
	"codeberg.org/reiver/go-log"
)

// Begin starts a new logging session and returns a Logger that writes
// structured log entries at the configured log level.
//
// Example usage:
//
//	log := logsrv.Begin()
//	defer log.End()
//	
//	log.Informf("server started")
func Begin(fields ...log.Field) log.Logger {
	return logger.Begin(fields...)
}
