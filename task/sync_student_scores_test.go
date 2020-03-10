package task

import (
	"testing"
	"time"
)

func TestUpdateStudentScore(t *testing.T) {
	_updateStudentScore()
	time.Sleep(5 * time.Second)
	t.Log("success")
}
