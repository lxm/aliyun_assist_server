package activation

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

type CreateActivationCodeReq struct {
	NamePrefix       string `form:"name_prefix" binding:"required"`
	ActiveCountLimit int    `form:"active_count_limit,default=-1"`
	Expire           int    `form:"expire,default=-1"`
	Description      string `form:"description"`
}

func CreateActivationCode(c *gin.Context) {
	var req CreateActivationCodeReq

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	activationCode, err := model.CreateActivationCode(req.NamePrefix, req.Description, req.ActiveCountLimit, req.Expire)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"activation_code": activationCode.Code,
	})

}
