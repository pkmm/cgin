package video_91porn

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func Test_video_91porn_parser(t *testing.T) {
	fmt.Println(len("中文"))
	result := ParseIndex()
	for _, item := range *result {
		log.Printf("%#v", item)
	}
}

func TestParseCategory(t *testing.T) {
	client := NewHttpClient()
	resq, _ := http.NewRequest("GET", "http://91porn.com/v.php?category=mf&viewtype=basic", nil)
	resp,_ := client.Do(resq)
	defer resp.Body.Close()
	result := parseByCategory(resp.Body)
	for _, item := range *result {
		log.Printf("%#v", item)
	}
}