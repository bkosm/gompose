package gompose

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDown(t *testing.T) {
	t.Run("down fails if there is no file", func(t *testing.T) {
		err := Down()
		assert.Error(t, err)
	})

	t.Run("down does not fail if there was no up before", func(t *testing.T) {
		err := Down(AsDownOpt(customFileOpt))
		assert.NoError(t, err)
	})

	t.Run("down cleans up after a successful setup", func(t *testing.T) {
		testUp(t)
		assertServiceIsUp(t)

		err := Down(AsDownOpt(customFileOpt))
		assert.NoError(t, err)

		assertServiceIsDown(t)
	})
}
