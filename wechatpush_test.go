package main

import (
	"cgin/service/wechatpush"
	"cgin/util"
	"testing"
)

func TestWxPusher(t *testing.T) {
	wp := wechatpush.NewPushBear([]string{"UID_Jo3zGWQgT9WmyKGQpVr4Oy2Juhkp"}, 1)
	worker := wechatpush.NewPusherContext(wp, "Hello", "pusher bear" + util.DateTime())
	_, err := worker.Push()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerSister(t *testing.T) {
	worker := wechatpush.NewPusherContext(wechatpush.NewPusherSister(), "hello", "server chan")
	_, err := worker.Push()
	if err != nil {
		t.Fatal(err)
	}
}