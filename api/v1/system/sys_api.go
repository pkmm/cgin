package system

import "github.com/gin-gonic/gin"

type SystemApi struct {
}

func (s *SystemApi) Index(c *gin.Context) {
	// TODO
	user, _ := apiService.GetUser("zcc")
	c.JSON(200, user)
}
