package gompose

import "errors"

type ReadyOrErrChan <-chan error

var ErrWaitTimedOut = errors.New("gompose: timed out waiting on condition")
