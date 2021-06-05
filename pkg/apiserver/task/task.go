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
		task := model.GetTaskByUUID(taskID)
		if task != nil {
			taskList.RunTasks = append(taskList.RunTasks, task.ParseRunTaskInfo())
		}
	}

	logrus.Infof("client list task with reason:%s", reason)
	c.JSON(200, taskList)
}
func Running(c *gin.Context) {
	taskID := c.Query("taskId")
	task := model.GetTaskByUUID(taskID)
	if task == nil {
		logrus.WithFields(logrus.Fields{
			"Module": "task",
			"Func":   "Finish",
		}).Errorf("Get task failed by taskId:%s", taskID)
		c.AbortWithStatus(200)
		return
	}

	rawData, _ := c.GetRawData()
	if len(rawData) > 0 {
		task.StashOutput(string(rawData))
	}
	c.AbortWithStatus(201)
}
func Finish(c *gin.Context) {
	//  "taskId" => "98",
	//  "start" => "1622876670110",
	//  "end" => "1622876679562",
	//  "exitCode" => "0",
	//  "dropped" => "0",
	taskID := c.Query("taskId")
	task := model.GetTaskByUUID(taskID)

	// start := c.Query("start")
	// end := c.Query("end")
	// exitCode := c.Query("exitCode")
	// dropped := c.Query("dropped")
	if task == nil {
		logrus.WithFields(logrus.Fields{
			"Module": "task",
			"Func":   "Finish",
		}).Errorf("Get task failed by taskId:%s", taskID)
		c.AbortWithStatus(200)
		return
	}

	rawData, _ := c.GetRawData()
	if len(rawData) > 0 {
		task.StashOutput(string(rawData))
		task.DumpOutput()
	}
	logrus.Infof("Task Finish:%v\n%v", string(rawData), c.Request.Header)
}
func Stopped(c *gin.Context) {}
func Timeout(c *gin.Context) {}
func Error(c *gin.Context)   {}
