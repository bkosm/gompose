# gompose
![Coverage](https://img.shields.io/badge/Coverage-97.8%25-brightgreen)
[![GoDoc](https://godoc.org/github.com/bkosm/gompose?status.svg)](https://godoc.org/github.com/bkosm/gompose)
[![CI](https://github.com/bkosm/gompose/actions/workflows/ci.yml/badge.svg)](https://github.com/bkosm/gompose/actions/workflows/ci.yml)
[![CodeQL](https://github.com/bkosm/gompose/actions/workflows/codeql.yml/badge.svg)](https://github.com/bkosm/gompose/actions/workflows/codeql.yml)
[![Go Report](https://goreportcard.com/badge/github.com/bkosm/gompose)](https://goreportcard.com/report/github.com/bkosm/gompose)

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

	g "github.com/bkosm/gompose/v2"
)

func TestMain(m *testing.M) {
	// Let's say we have a postgres container in the spec
	err := g.Up(
		g.Wait(
			g.ReadyOnLog("database system is ready to accept connections", g.Times(2)),
		),
		g.SignalCallback(func(_ os.Signal) { // any action to be taken on SIGINT, SIGTERM
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

#### but I have things to do in the meantime

Waiting is done the idiomatic way - you can await the ready channel without passing it to `Up` whenever:

```go
//...
err := g.Up(g.SignalCallback(func(_ os.Signal) { _ = g.Down() }))
if err != nil {
	log.Fatal(err)
}

// do stuff

readyOrErr := g.ReadyOnLog("database system is ready to accept connections", g.Times(2))
if err := <-readyOrErr; err != nil {
	log.Fatal(err)
}

code := m.Run()
//...
```

#### but I want to wait until service passes health-checks

This can be done by using the `ReadyOnHttp` wait channel:

```go
hc, _ := http.NewRequest(http.MethodGet, "http://localhost:5432", nil)

err := g.Up(g.Wait(g.ReadyOnHttp(*hc)))
```

And you can customize what it means to be healthy too:

```go
err := g.Up(g.Wait(
    g.ReadyOnHttp(
        *hc,
        g.ResponseVerifier(func (r *http.Response) (bool, error) {
            return r.StatusCode == http.StatusUnauthorized, nil
        })
    ),
))
```

#### miscellaneous

##### skip down non-programmatically

When an environment variable of `GOMPOSE_SKIP` is set with a value including `down`
(as a comma delimited string), the invocation of `docker-compose down` is skipped, leaving out the
environment in the state it was during the tests. This can be helpful while debugging and
during development, in cases where the services needed require significant effort to start.

##### all the commands with all their options - as examples

```go
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	g "github.com/bkosm/gompose"
)

func main() {
	// ready on stdout
	cmd := exec.Command("echo", "1")
	wait := g.Wait(
		g.ReadyOnStdout(
			cmd,
			"1",
			g.Times(1),
			g.Timeout(time.Minute),
			g.PollInterval(time.Second),
		),
	)

	// ready on log
	wait = g.Wait(
		g.ReadyOnLog(
			"log line",
			g.CustomFile("other/compose.yaml"),
			g.Times(3),
			g.Timeout(2*time.Minute),
			g.PollInterval(3*time.Second),
		),
	)

	// ready on http
	payload := `{
            "is_fine": true
        }`
	req, _ := http.NewRequest(
		http.MethodPut,
		"http://myservice/api",
		bytes.NewReader([]byte(payload)),
	)
	wait = g.Wait(
		g.ReadyOnHttp(
			*req,
			g.Times(1),
			g.Timeout(time.Minute),
			g.PollInterval(time.Second),
			g.ResponseVerifier(func(_ *http.Response) (bool, error) {
				return true, nil
			}),
		),
	)

	// Up
	_ = g.Up(
		g.CustomFile("other/compose.yaml"),
		g.CustomServices("s1", "s2"),
		g.SignalCallback(func(s os.Signal) {
			fmt.Println("interrupt!")
		}),
		wait,
	)

	// Down
	_ = g.Down(
		g.CustomFile("other/compose.yaml"),
	)
}
```

## contributing

This is the absolute bare-bone of a library and contributions are welcome.

List of things to do in priority:

1. Covering a larger portion of CLI's capabilities
1. More wait conditions

When contributing, be mindful that the purpose of this library is to be
self-contained, lean and easy to use.

## licensing

This software is open-source. See [license](LICENSE).
