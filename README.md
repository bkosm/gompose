# gompose

Straightforward library to use `docker-compose` programmatically in Go tests.

Unlike other (deprecated) libraries, this one relies on `os/exec` instead of
heavy docker libraries and others.

It's just enough so that you can set up and clean your environment for some automated tests.

The only hard requirement is `docker-compose` executable on system path.

## example

```go
package main

import (
	"log"
	"os"
	"testing"

	"github.com/bkosm/gompose"
)

func TestMain(m *testing.M) {
	var err error

	// We have a postgres container in the spec
	err = gompose.Up(
		gompose.ReadyOnLog("database system is ready to accept connections", 2),
		func() { _ = gompose.Down() }, // any action to be taken on system interrupt
	)
	if err != nil {
		log.Panic(err)
	}

	code := m.Run()

	err = gompose.Down()
	if err != nil {
		log.Panic(err)
	}

	os.Exit(code)
}
```

When you run `go test ./...` you will see that the container is starting and running before any tests.
After the tests conclude, `compose down` is performed.

In case of a system interrupt (`SIGINT`, `SIGTERM`), the library allows for custom callbacks to ensure that no dangling
containers are left after the user requested a stop.

## contributing

This is the absolute bare-bone of a library and contributions are welcome.

The first thing on the list are some tests,
then documentation,
then more configurability and
then covering a larger portion of the cli's capabilities and wait conditions.

When contributing, be mindful that the purpose of this library is to be
self-contained, lean and easy to use.

## licensing

This software is open-source. See [license](LICENSE).
