package zf

import "testing"

func TestZfGetScores(t *testing.T) {
	crawl, err := NewCrawl("201623200101025", "200011da")

	if err != nil {
		t.Error(err)
		return
	}
	ret, err := crawl.GetScores()
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 0 {
		//t.Error("成绩列表不为空")
		t.Error(ret)
		return
	}
}
