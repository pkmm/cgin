package main

import (
	"cgin/service"
	"testing"
)

func TestWeChatNotificationService_Notify(t *testing.T) {
	_, err := service.WeChatNotificationService.Notify(
		"This is title",
		"This is content.",
		service.ServerSister,
		service.PlainText,
		[]string{""},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.WeChatNotificationService.Notify(
		"This is title",
		"This is content.",
		service.WxPusher,
		service.PlainText,
		[]string{"UID_Jo3zGWQgT9WmyKGQpVr4Oy2Juhkp"},
	)
	if err != nil {
		t.Fatal(err)
	}
}
