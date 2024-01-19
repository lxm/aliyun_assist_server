package model

import (
	"errors"
	"fmt"

	"github.com/lxm/aliyun_assist_server/pkg/util"
)

type ActivationCode struct {
	ID               int    `json:"id" gorm:"primarykey"`
	Code             string `json:"code" gorm:"type:varchar(100);uniqueIndex"`
	InstanceID       string `json:"instance_id" gorm:"type:varchar(50);index"`
	NamePrefix       string `json:"name_prefix" gorm:"varchar(50);uniqueIndex;size:50"`
	ActiveCountLimit int    `json:"active_count_limit"`
	ActiveCountUsed  int    `json:"active_count_used"`
	Expire           int    `json:"expire"`
	Description      string `json:"description" gorm:"type:varchar(200)"`
}

func CreateActivationCode(namePrefix, description string, activeCountLimit, expire int) (*ActivationCode, error) {
	code := "a-" + util.RandStringRunes(10)

	ac := &ActivationCode{
		Code:             code,
		NamePrefix:       namePrefix,
		ActiveCountLimit: activeCountLimit,
		Expire:           expire,
		Description:      description,
	}
	if activeCountLimit == 1 {
		ac.InstanceID = "i-" + util.RandStringRunes(16)
	}
	err := db.Create(&ac).Error
	return ac, err
}

func CheckActivationCode(code string) (*ActivationCode, string, error) {
	var ac ActivationCode

	err := db.Where("code = ?", code).First(&ac).Error
	if err != nil {
		return nil, "", err
	}

	if ac.ActiveCountLimit == -1 || ac.ActiveCountUsed < ac.ActiveCountLimit {
		ac.ActiveCountUsed = ac.ActiveCountUsed + 1
		name := fmt.Sprintf("%s-%d", ac.NamePrefix, ac.ActiveCountUsed)
		err := db.Save(ac).Error
		if err != nil {
			return nil, "", err
		}
		return &ac, name, err
	} else {
		return nil, "", errors.New("code used over limit")
	}
}
