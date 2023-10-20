package gompose

import (
	"testing"
)

func TestAliases(t *testing.T) {
	t.Run("PostgresViaLogs", func(t *testing.T) {
		_ = PostgresViaLogs()
	})

	t.Run("DownOnSignal", func(t *testing.T) {
		_ = DownOnSignal(customFileOpt)
	})
}
