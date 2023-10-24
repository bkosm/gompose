package gompose

import (
	"os"
	"time"
)

// PostgresViaLogs is an alias for a ReadyOrErrChan that checks if a Postgres
// instance is ready by checking for a characteristic log entry.
func PostgresViaLogs(opts ...Option) ReadyOrErrChan {
	options := append([]Option{Times(2)}, opts...)
	return ReadyOnLog("database system is ready to accept connections", options...)
}

// DownOnSignal is an alias for an option that provides a SignalCallback which performs a Down in
// case of any system interrupt.
// Make sure to provide the CustomFile as a parameter if such is used for Up.
func DownOnSignal(opts ...Option) Option {
	return SignalCallback(func(_ os.Signal) { _ = Down(opts...) })
}

// Retry is an alias for a RetryCommand with sensible defaults.
// It will cause the command to run 3 times with 2 seconds in between before returning an error.
func Retry() Option {
	return RetryCommand(3, time.Second*2)
}
