package main

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestRun(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	testFile := "testdata/noprogramchange.js"

	c.Assert(run(testFile), qt.IsNil)
}
