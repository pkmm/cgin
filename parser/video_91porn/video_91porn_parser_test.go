package video_91porn

import (
	"log"
	"testing"
)

func Test_video_91porn_parser(t *testing.T) {
	result := ParseIndex()
	for _, item := range *result {
		log.Printf("%#v", item)
	}
}