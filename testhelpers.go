package gompose

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os/exec"
	"testing"
)

const (
	expectedLine      = "server is listening"
	expectedResponse  = "ok"
	customServiceName = "echo"
	containerPort     = 5678
)

func testUp(t *testing.T) {
	_, err := run(*exec.Command("docker-compose", "-f", "./testdata/docker-compose.yml", "up", "-d"))
	assert.NoError(t, err)
}

func testDown(t *testing.T) {
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
