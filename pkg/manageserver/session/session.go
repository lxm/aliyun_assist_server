package session

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/lxm/aliyun_assist_server/pkg/util"
)

type SessionTaskInfo struct {
	CmdContent   string `json:"cmdContent"`
	Username     string `json:"username"`
	Password     string `json:"windowsPasswordName"`
	SessionId    string `json:"channelId"`
	WebsocketUrl string `json:"websocketUrl"`
	TargetHost   string `json:"targetHost"`
	PortNumber   string `json:"portNumber"`
	FlowLimit    int    `json:"flowLimit"` // 最大流量 单位 bps
}

func StartSession(c *gin.Context) {
	instanceID := c.Query("instance_id")
	if instanceID == "" {
		c.JSON(400, gin.H{
			"message": "instance_id is required",
		})
		return
	}
	invokeId := "t-" + util.RandStringRunes(16)

	redisClient := redisclient.GetClient()

	ctx := context.Background()
	channel := "notify_server:" + instanceID
	msg := fmt.Sprintf("kick_vm session start %s", invokeId)
	_, err := redisClient.Publish(ctx, channel, msg).Result()

	c.JSON(200, gin.H{
		"invoke_id": invokeId,
		"err":       err,
	})

}
