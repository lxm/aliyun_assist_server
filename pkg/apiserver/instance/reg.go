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
		code := regInfo.ActivationCode
		// logrus.Infof("code:%v", code)
		// get exist reginfo
		existRegInfos := model.GetRegisterInfoByActivactionCode(code)
		if len(*existRegInfos) == 1 {
			existRegInfo := (*existRegInfos)[0]
			existRegInfo.PublicKeyBase64 = regInfo.PublicKeyBase64
			existRegInfo.MachineId = regInfo.MachineId
			existRegInfo.Hostname = regInfo.Hostname
			existRegInfo.InstanceName = regInfo.InstanceName
			existRegInfo.IntranetIp = regInfo.IntranetIp
			model.GetDB().Model(&model.RegisterInfo{}).Where("id = ?", existRegInfo.ID).Updates(existRegInfo)
			c.JSON(200, gin.H{
				"code":       200,
				"instanceId": existRegInfo.InstanceID,
			})
			return
		}
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
		if ac.ActiveCountLimit == 1 && len(ac.InstanceID) > 0 {
			regInfo.InstanceID = ac.InstanceID
		} else {
			regInfo.GenInstanceID()
		}
		regInfo.InstanceName = instanceName
		model.GetDB().Save(&regInfo)

		c.JSON(200, gin.H{
			"code":       200,
			"instanceId": regInfo.InstanceID,
		})
		return
	}
	// var registerResponse map[string]interface{}

}
func DeReg(c *gin.Context) {}
