package util

import (
	"path"
	"testing"
)

func TestGetCurrentCodePath(t *testing.T) {
	t.Log(GetExecutableCodePath())
	t.Log(GetCurrentCodePath())
	t.Log(path.Join(GetCurrentCodePath(), ".."))

	t.Log(TimeString())
}
