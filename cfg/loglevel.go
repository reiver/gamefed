package cfg

import (
	"codeberg.org/reiver/go-env"
)

// LogLevel returns the configured log verbosity level, read from the
// LOG_LEVEL environment variable. If the variable is unset or invalid,
// the default level "vvv" is returned. Valid values are "v" through
// "vvvvvvv".
//
// Example usage:
//
//	var level string = cfg.LogLevel()
func LogLevel() string {
	const defaultLogLevel string = "vvv"

	var value string = env.GetElse[string]("LOG_LEVEL", defaultLogLevel)
	switch value {
	case "v","vv","vvv","vvvv","vvvvv","vvvvvv","vvvvvvv":
		// nothing here
	default:
		value = defaultLogLevel
	}

	return value
}
