package system

import "github.com/gin-gonic/gin"

type SystemApi struct {
}

func (s *SystemApi) Index(c *gin.Context) {
	// TODO
	user, _ := apiService.GetUser("zcc")
	_, t := deliAutoSignService.SignOne(user)
	c.JSON(200, t)
}
