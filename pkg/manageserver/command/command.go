package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
)

type RunCommandReq struct {
	Command model.Command `json:"command"`
	Options struct {
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
	task := model.CreateTask(runCmdReq.Command.ID, runCmdReq.Options.TaskOption)

	redisClient := redisclient.GetClient()

	ctx := context.Background()
	instanceIDList := strings.Split(runCmdReq.Options.InstanceIDs, ",")
	for _, instanceID := range instanceIDList {
		channel := "notify_server:" + instanceID
		msg := fmt.Sprintf("kick_vm task run %s", task.UUID)
		redisClient.Publish(ctx, channel, msg)
	}

	c.JSON(200, runCmdReq)

}
