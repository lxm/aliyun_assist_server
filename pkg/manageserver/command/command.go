package command

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/lxm/aliyun_assist_server/pkg/util"
)

type RunCommandReq struct {
	Command     model.Command `json:"command" binding:"required"`
	InstanceIDs []string      `json:"instance_ids" binding:"required"`
	Options     struct {
		model.TaskOption
		KeepCommand     bool   `json:"keep_command"`
		ContentEncoding string `json:"content_encoding" binding:"required"`
	} `json:"options" binding:"required"`
}

func RunCommand(c *gin.Context) {
	var runCmdReq RunCommandReq
	err := c.ShouldBindJSON(&runCmdReq)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = model.CreateCommand(&runCmdReq.Command)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	invokeId := "t-" + util.RandStringRunes(16)

	redisClient := redisclient.GetClient()

	ctx := context.Background()
	for _, instanceID := range runCmdReq.InstanceIDs {
		task := model.CreateTask(runCmdReq.Command.ID, instanceID, invokeId, runCmdReq.Options.TaskOption)
		channel := "notify_server:" + instanceID
		msg := task.KickMsg()
		num, err := redisClient.Publish(ctx, channel, msg).Result()
		if err == nil && num == 0 {
			redisClient.SAdd(ctx, channel+"no_sub", msg)
		}
	}

	c.JSON(200, gin.H{
		"command_id": runCmdReq.Command.ID,
		"invoke_id":  invokeId,
	})

}

func InvokeCommand(c *gin.Context) {
}
