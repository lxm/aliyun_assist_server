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
		"instance_id":     activationCode.InstanceID,
	})

}

func GetInstancesByActivationCode(c *gin.Context) {
	code := c.Param("code")
	instances := model.GetRegisterInfoByActivactionCode(code)

	c.JSON(200, gin.H{
		"instances": instances,
	})
}

func BatchGetInstancesByActivationCode(c *gin.Context) {
	codes := c.QueryArray("codes")
	instances := model.BatchGetRegisterInfoByActivactionCode(codes)

	instanceMap := make(map[string]model.RegisterInfo)

	for _, v := range *instances {
		instanceMap[v.ActivationCode] = v
	}

	c.JSON(200, gin.H{
		"instances": instanceMap,
	})
}
