package require

import (
	"testing"

	"github.com/alexandreh2ag/htransformation/pkg/tests/assert"
)

func NoError(t *testing.T, err error) {
	t.Helper()

	if !assert.NoError(t, err) {
		t.FailNow()
	}
}
