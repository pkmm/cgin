package task

import "testing"

func TestUpdateStudentScore(t *testing.T) {
	_updateStudentScore()
	t.Log("success")
	select {} // 阻塞使用线程池
}
