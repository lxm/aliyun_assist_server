package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
)

const (
	INSTANCE_STATUS_HASH = "instance:status"
	STATUS_ONLINE        = 1
	STATUS_OFFLINE       = 0
	MAX_STATUS_TTL       = 600
)

type InstanceStatus struct {
	UpdateTime int64 `json:"update_time"`
	Status     int64 `json:"status"`
}
type Instance struct {
	Status string `json:"status"`
	*RegisterInfo
}

func GetIntanceByID(ID string) *Instance {
	instance := &Instance{}
	regInfo := GetRegisterInfo(ID)
	if regInfo == nil {
		return nil
	}
	instance.RegisterInfo = regInfo
	return instance
}

func GetInstanceList() {

}

func (i *Instance) GetStatus() int64 {
	ctx := context.Background()
	redisClient := redisclient.GetClient()

	statusStr, err := redisClient.HGet(ctx, INSTANCE_STATUS_HASH, i.InstanceID).Result()
	if err != nil {
		return STATUS_OFFLINE
	}
	var status InstanceStatus

	err = json.Unmarshal([]byte(statusStr), &status)
	if err != nil {
		return STATUS_OFFLINE
	}
	if time.Now().Unix()-status.UpdateTime > MAX_STATUS_TTL {
		return STATUS_OFFLINE
	}
	return status.Status
}

func (i *Instance) SetOnline() error {
	status := &InstanceStatus{
		UpdateTime: time.Now().Unix(),
		Status:     1,
	}
	ctx := context.Background()
	statusStr, err := json.Marshal(status)
	if err != nil {
		return err
	}
	redisClient := redisclient.GetClient()
	return redisClient.HSet(ctx, INSTANCE_STATUS_HASH, i.InstanceID, string(statusStr)).Err()
}

func (i *Instance) SetOffline() error {
	status := &InstanceStatus{
		UpdateTime: time.Now().Unix(),
		Status:     0,
	}
	ctx := context.Background()
	statusStr, err := json.Marshal(status)
	if err != nil {
		return err
	}
	redisClient := redisclient.GetClient()
	return redisClient.HSet(ctx, INSTANCE_STATUS_HASH, i.InstanceID, string(statusStr)).Err()
}
