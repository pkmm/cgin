package service

import "testing"

func TestWeChatNotificationService_Notify(t *testing.T) {
	_, err := WeChatNotificationService.Notify(
		"This is title",
		"This is content.",
		ServerSister,
		PlainText,
		[]string{""},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = WeChatNotificationService.Notify(
		"This is title",
		"This is content.",
		WxPusher,
		PlainText,
		[]string{"UID_Jo3zGWQgT9WmyKGQpVr4Oy2Juhkp"},
	)
	if err != nil {
		t.Fatal(err)
	}
}
