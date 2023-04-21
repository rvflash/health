package health_test

import (
	"context"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/health"
)

func TestCreateFileCheck(t *testing.T) {
	t.Parallel()

	is.New(t).NoErr(health.CreateFileCheck(os.TempDir(), "")(context.Background()))
}
