package signalerror

import (
	"errors"
	"os"
	"syscall"
)

// NewSignalError returns a new error wrapping a signal
func NewSignalError(signal os.Signal) *ErrSignal {
	return &ErrSignal{signal}
}

// ErrSignal is returned as error when a signal was caught
type ErrSignal struct {
	Signal os.Signal
}

func (e *ErrSignal) Error() string {
	return "caught signal " + e.Signal.String()
}

// ExitCode returns the best-current-practice exit code for a signal
func (e *ErrSignal) ExitCode() int {
	origSignal, ok := e.Signal.(syscall.Signal)
	if !ok {
		origSignal = 0
	}
	return 128 + int(origSignal)
}

// ErrSignalExitCode returns the signal's best-current-practice exit code, if err was a ErrSignal
func ErrSignalExitCode(err error) (int, bool) {
	var errSignal *ErrSignal
	if errors.As(err, &errSignal) {
		return errSignal.ExitCode(), true
	}
	return 0, false
}

// ErrSignalIsTerm returns true if err was a ErrSignal for SIGTERM
func ErrSignalIsTerm(err error) bool {
	var errSignal *ErrSignal
	if errors.As(err, &errSignal) {
		return errSignal.Signal == syscall.SIGTERM
	}
	return false
}
