package task

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/types"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/sirupsen/logrus"
)

func Invalid(c *gin.Context) {}
func List(c *gin.Context) {
	// instanceID := c.GetString("checked-instance-id")
	var taskList types.TaskListResp
	taskList.RunTasks = make([]interface{}, 0)
	// taskList.StopTasks = []model.RunTaskInfo{}
	// taskList.SendFileTasks = []model.SendFileTaskInfo{}

	reason := c.Query("reason")
	taskID := c.Query("taskId")
	logrus.Infof("taskId:%v", taskID)
	if reason == "kickoff" {
		task := model.GetTaskByID(taskID)
		if task != nil {
			taskList.RunTasks = append(taskList.RunTasks, task.ParseRunTaskInfo())
		}
	}

	logrus.Infof("client list task with reason:%s", reason)
	c.JSON(200, taskList)
}
func Running(c *gin.Context) {}
func Finish(c *gin.Context)  {}
func Stopped(c *gin.Context) {}
func Timeout(c *gin.Context) {}
func Error(c *gin.Context)   {}
