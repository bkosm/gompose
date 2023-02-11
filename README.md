# gompose
[![](https://github.com/bkosm/gompose/actions/workflows/ci.yml/badge.svg)](https://github.com/bkosm/gompose/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-96.3%25-brightgreen)

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

	g "github.com/bkosm/gompose"
)

func TestMain(m *testing.M) {
	// Let's say we have a postgres container in the spec
	err := g.Up(
		g.WithWait(
			g.ReadyOnLog(g.WithText("database system is ready to accept connections"), g.Times(2)),
		),
		g.WithSignalCallback(func(_ os.Signal) { // any action to be taken on SIGINT, SIGTERM
			_ = g.Down()
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err = g.Down(); err != nil {
		log.Fatal(err)
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

List of things to do in priority:

1. ~~Tests~~
1. CI and badges
1. Documentation
1. ~~Configurability~~
1. Covering a larger portion of CLI's capabilities
1. More wait conditions

When contributing, be mindful that the purpose of this library is to be
self-contained, lean and easy to use.

## licensing

This software is open-source. See [license](LICENSE).
