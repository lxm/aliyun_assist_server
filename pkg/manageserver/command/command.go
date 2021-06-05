package command

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
)

type RunCommandReq struct {
	Command     model.Command `json:"command"`
	InstanceIDs []string      `json:"instance_ids"`
	Options     struct {
		model.TaskOption
		KeepCommand     bool   `json:"keep_command"`
		ContentEncoding string `json:"content_encoding"`
	} `json:"options"`
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
	// invokeId := util.RandStringRunes(32)

	redisClient := redisclient.GetClient()

	ctx := context.Background()
	for _, instanceID := range runCmdReq.InstanceIDs {
		task := model.CreateTask(runCmdReq.Command.ID, instanceID, runCmdReq.Options.TaskOption)
		channel := "notify_server:" + instanceID
		msg := fmt.Sprintf("kick_vm task run %s", task.UUID)
		redisClient.Publish(ctx, channel, msg)
	}

	c.JSON(200, runCmdReq)

}
