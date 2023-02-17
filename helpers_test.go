package gompose

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

const (
	expectedLogLine   = "server is listening"
	customServiceName = "echo"
	containerPort     = 5678
)

var customFileOpt = WithCustomFile("./testdata/docker-compose.yml")

func testUp(t *testing.T) {
	t.Helper()
	_, err := run(*exec.Command("docker-compose", "-f", "./testdata/docker-compose.yml", "up", "-d"))
	if err != nil {
		t.Fatal(err)
	}
}

func testDown(t *testing.T) {
	t.Helper()
	_, err := run(*exec.Command("docker-compose", "-f", "./testdata/docker-compose.yml", "down"))
	if err != nil {
		t.Fatal(err)
	}
}

func pingService() error {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", containerPort))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("ping status code was not 200")
	}
	return nil
}

func serviceIsUp() bool {
	err := pingService()
	return err == nil
}

func assertServiceIsUp(t *testing.T) {
	t.Helper()
	err := pingService()
	if err != nil {
		t.Fatal("expected service to be up, got", err)
	}
}

func serviceIsDown() bool {
	err := pingService()
	return err != nil
}

func assertServiceIsDown(t *testing.T) {
	t.Helper()
	err := pingService()
	if err == nil {
		t.Fatal("service is up")
	}
}

func goIntoTestDataDir(t *testing.T) func() {
	t.Helper()
	startDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(fmt.Sprintf("%s/testdata", startDir))
	if err != nil {
		t.Fatal(err)
	}

	return func() {
		err = os.Chdir(startDir)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func doSignal(t *testing.T, s syscall.Signal) {
	t.Helper()
	err := syscall.Kill(syscall.Getpid(), s)
	if err != nil {
		t.Fatal(err)
	}
}

func validRequest(t *testing.T) *http.Request {
	t.Helper()

	return MustT[*http.Request](t)(http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:%d", containerPort),
		nil,
	))
}

// assertEventually is a helper function copied from stretchr/testify
func assertEventually(t *testing.T, condition func() bool, waitFor time.Duration, tick time.Duration) {
	t.Helper()

	ch := make(chan bool, 1)

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			t.Fatal("eventual condition not satisfied")
			return
		case <-tick:
			tick = nil
			go func() { ch <- condition() }()
		case v := <-ch:
			if v {
				return
			}
			tick = ticker.C
		}
	}
}

func assertError(t *testing.T, errs ...error) {
	t.Helper()
	l := len(errs)
	switch l {
	case 0:
		return
	case 1:
		if errs[0] == nil {
			t.Fatal("expected error, got nil")
		}
	default:
		expected := errs[l-1]
		for _, err := range errs[:l-1] {
			if err == nil || !errors.Is(err, expected) {
				t.Fatalf("expected error of %v but got %v instead", expected, err)
			}
		}
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("expected no error, got", err)
	}
}
