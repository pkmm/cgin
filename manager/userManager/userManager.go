package userManager

import (
	"pkmm_gin/model"
	"pkmm_gin/utility"
	"strconv"
	"time"
)

func GenerateUserToken(user *model.User, openId string) (token string) {
	token, _ = utility.GenerateSignatureAndId(map[string]string{
		"open_id": openId,
		"user_id": strconv.FormatInt(user.Id, 10),
	})

	myRedis := model.GetRedis()
	defer myRedis.Close()

	// 最新的数据 token
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	myRedis.Do("SET", user.Id, token)
	myRedis.Do("EXPIRE", user.Id, expireAt)

	return token
}

