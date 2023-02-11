package gompose

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUp(t *testing.T) {
	t.Run("up works with no arguments", func(t *testing.T) {

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
	})

	t.Run("intercepts os signals", func(t *testing.T) {

	})
}
