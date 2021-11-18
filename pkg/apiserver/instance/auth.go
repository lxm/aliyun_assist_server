package instance

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/util"
	"github.com/sirupsen/logrus"
)

func CheckHeaderMiddleware(c *gin.Context) {
	instanceID := c.GetHeader("x-acs-instance-id")
	timestamp := c.GetHeader("x-acs-timestamp")
	requestID := c.GetHeader("x-acs-request-id")
	signature := c.GetHeader("x-acs-signature")
	regInfo := model.GetRegisterInfo(instanceID)
	if regInfo != nil && regInfo.ID != 0 {
		mid := regInfo.MachineId
		publicKeyBase64 := regInfo.PublicKeyBase64
		publicKeyByte, err := base64.StdEncoding.DecodeString(publicKeyBase64)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func":    "CheckHeaderMiddleware",
				"section": "SignDecodeString",
			}).Errorf("SignDecodePublicKey err :%v", err)
		}
		input := instanceID + mid + timestamp + requestID
		sign, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func":    "CheckHeaderMiddleware",
				"section": "SignDecodeString",
			}).Errorf("SignDecodeString err :%v", err)
		}
		err = util.RsaCheckSign(input, sign, string(publicKeyByte))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func":    "CheckHeaderMiddleware",
				"section": "RsaCheckSign",
			}).Errorf("RsaCheckSign err :%v", err)
		} else {
			c.Set("checked-instance-id", instanceID)
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"func": "CheckHeaderMiddleware",
		}).Errorf("GetRegisterInfo with instanceID:%v failed", instanceID)
	}
}

func checkSign() {

}

// header.Add("x-acs-instance-id", instance_id)
// header.Add("x-acs-timestamp", str_timestamp)
// header.Add("x-acs-request-id", str_request_id)
// header.Add("x-acs-signature", output)
