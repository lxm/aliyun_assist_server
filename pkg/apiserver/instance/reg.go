package instance

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/sirupsen/logrus"
)

func Reg(c *gin.Context) {
	rawData, _ := c.GetRawData()
	logrus.Infof("Reg RawData:%v", string(rawData))
	var regInfo model.RegisterInfo

	err := json.Unmarshal(rawData, &regInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":       400,
			"instanceId": "",
		})
		return
	} else {
		instanceID := regInfo.GenInstanceID()
		code := regInfo.ActivationCode
		logrus.Infof("code:%v", code)
		ac, instanceName, err := model.CheckActivationCode(code)
		if ac == nil {
			c.JSON(400, gin.H{
				"message": "invaild active code",
			})
			return
		}
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		regInfo.InstanceName = instanceName
		model.GetDB().Save(&regInfo)

		c.JSON(200, gin.H{
			"code":       200,
			"instanceId": instanceID,
		})
		return
	}
	// var registerResponse map[string]interface{}

}
func DeReg(c *gin.Context) {}
