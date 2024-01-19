package model

import (
	"github.com/lxm/aliyun_assist_server/pkg/util"
	"gorm.io/gorm"
)

type RegisterInfo struct {
	gorm.Model
	ActivationCode  string `json:"activationCode" gorm:"type:varchar(256);index"`
	MachineId       string `json:"machineId" gorm:"type:varchar(256);index"`
	RegionId        string `json:"regionId" gorm:"type:varchar(100)"`
	InstanceName    string `json:"instanceName" gorm:"type:varchar(256);index"`
	Hostname        string `json:"hostname" gorm:"type:varchar(256)"`
	IntranetIp      string `json:"intranetIp" gorm:"type:varchar(256)"`
	OsVersion       string `json:"osVersion" gorm:"type:varchar(256)"`
	OsType          string `json:"osType" gorm:"type:varchar(256)"`
	ClientVersion   string `json:"agentVersion" gorm:"type:varchar(256)"`
	PublicKeyBase64 string `json:"publicKey" gorm:"type:varchar(3000)"`
	InstanceID      string `json:"InstanceId" gorm:"type:varchar(256);index"`
	ActivationID    string `json:"activationId" bgorm:"type:varchar(256);index"`
}

func GetRegisterInfo(instanceID string) *RegisterInfo {
	var regInfo RegisterInfo
	err := GetDB().Where("instance_id", instanceID).Find(&regInfo).Error
	if err != nil {
		return nil
	}
	return &regInfo
}

func GetRegisterInfoByActivactionCode(code string) *[]RegisterInfo {
	var regInfo []RegisterInfo
	err := GetDB().Where("activation_code", code).Find(&regInfo).Error
	if err != nil {
		return nil
	}
	return &regInfo
}

func (r *RegisterInfo) GenInstanceID() string {
	r.InstanceID = "i-" + util.RandStringRunes(16)
	return r.InstanceID
}

func BatchGetRegisterInfoByActivactionCode(codes []string) *[]RegisterInfo {
	var regInfo []RegisterInfo
	err := GetDB().Where("activation_code in (?)", codes).Find(&regInfo).Error
	if err != nil {
		return nil
	}
	return &regInfo
}
