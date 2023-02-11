package gompose

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

const (
	expectedLine      = "Server is listening"
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
