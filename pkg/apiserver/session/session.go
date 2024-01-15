package session

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/types"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func List(c *gin.Context) {
	channelID := c.Query("channelId")
	if channelID == "" {
		c.JSON(400, gin.H{
			"message": "channelId is required",
		})
		return
	}
	// get task by taskUUID

	var taskList types.TaskListResp
	taskList.SessionTasks = make([]interface{}, 0)
	task := model.GetTaskByUUID(channelID)
	if task != nil {
		taskList.SessionTasks = append(taskList.SessionTasks, task.ParseSessionTaskInfo())
	}
	c.JSON(200, taskList)
}

func Websocket(c *gin.Context) {

}
