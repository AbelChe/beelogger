package beelogger

import (
	"fmt"
	"testing"
)

func TestBeeLogger(t *testing.T) {
	Info().Msg("This is an info message")
	Warn().Msg("This is a warning message")
	Error().Msg("This is an error message")
	Debug().Msg("This is an debug message")
	Err(fmt.Errorf("this is an error object for %s", "testing")).Msg("This is an error message with error object")
}
