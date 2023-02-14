package gompose

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
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
	assert.NoError(t, err)
}

func testDown(t *testing.T) {
	t.Helper()
	_, err := run(*exec.Command("docker-compose", "-f", "./testdata/docker-compose.yml", "down"))
	require.NoError(t, err)
}

func pingService() error {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", containerPort))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}

func assertServiceIsUp(t *testing.T) {
	t.Helper()
	err := pingService()
	assert.NoError(t, err)
}

func serviceIsDown() bool {
	err := pingService()
	return err != nil
}

func assertServiceIsDown(t *testing.T) {
	t.Helper()
	err := pingService()
	assert.Error(t, err)
}

func goIntoTestDataDir(t *testing.T) func() {
	t.Helper()
	startDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(fmt.Sprintf("%s/testdata", startDir))
	require.NoError(t, err)

	return func() {
		err = os.Chdir(startDir)
		require.NoError(t, err)
	}
}

func doSignal(t *testing.T, s syscall.Signal) {
	t.Helper()
	err := syscall.Kill(syscall.Getpid(), s)
	require.NoError(t, err)
}

func validRequest(t *testing.T) *http.Request {
	t.Helper()
	rq, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:%d", containerPort),
		nil,
	)
	assert.NoError(t, err)
	return rq
}
