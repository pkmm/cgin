package util

import "math/rand"

func RandomString(length int) string {
	var letters = []rune("abcdefg-*%$#@=&^hijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ9876543210")
	N := len(letters)
	ans := make([]rune, length)
	for i := 0; i < length; i++ {
		ans[i] = letters[rand.Intn(N)]
	}
	return string(ans)
}
