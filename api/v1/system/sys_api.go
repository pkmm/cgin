package system

import (
	"cgin/global"
	"github.com/gin-gonic/gin"
)

type SystemApi struct {
}

func (s *SystemApi) Index(c *gin.Context) {
	// TODO
	t := global.Config.Deli.Season
	c.JSON(200, t)
}
