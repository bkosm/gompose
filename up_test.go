package gompose

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUp(t *testing.T) {
	t.Run("up works with no arguments", func(t *testing.T) {
		goBack := goIntoTestDataDir(t)
		defer func() {
			goBack()
			testDown(t)
		}()

		err := Up()
		assert.NoError(t, err)
		assertServiceIsUp(t)
	})

	customFileOpt := WithCustomFile("./testdata/docker-compose.yml")

	t.Run("up works with options", func(t *testing.T) {
		defer testDown(t)

		err := Up(
			WaitFor(ReadyOnLog(Text(expectedLine), AsReadyOpt(customFileOpt))),
			WithCustomServices(customServiceName),
			AsUpOpt(customFileOpt),
		)
		assert.NoError(t, err)
		assertServiceIsUp(t)
	})

	t.Run("intercepts os signals", func(t *testing.T) {

	})
}

func assertServiceIsUp(t *testing.T) {
	err := pingService()
	assert.NoError(t, err)
}

func goIntoTestDataDir(t *testing.T) func() {
	startDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(fmt.Sprintf("%s/testdata", startDir))
	require.NoError(t, err)

	return func() {
		err = os.Chdir(startDir)
		require.NoError(t, err)
	}
}
