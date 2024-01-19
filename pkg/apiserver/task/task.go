package task

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/types"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/sirupsen/logrus"
)

const (
	FetchOnKickoff string = "kickoff"
	FetchOnStartup string = "startup"
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

	//debug
	logrus.Infof("taskId:%v reason:%v", taskID, reason)

	if reason == FetchOnKickoff {
		task := model.GetTaskByUUID(taskID)
		if task != nil {
			taskList.RunTasks = append(taskList.RunTasks, task.ParseRunTaskInfo())
		}
	} else if reason == FetchOnStartup {
		defer func() {
			instanceID := c.GetHeader("x-acs-instance-id")
			ctx := context.Background()
			redisClient := redisclient.GetClient()
			channel := "notify_server:" + instanceID
			msgs, _ := redisClient.SMembers(ctx, channel+"no_sub").Result()
			for _, msg := range msgs {
				redisClient.Publish(ctx, channel, msg).Result()
			}
			redisClient.Del(ctx, channel+"no_sub")
		}()
	}

	if reason == "startup" {
		tasks := model.GetTasksByStatus(model.TASK_STATUS_PENDING)
		for _, task := range tasks {
			taskList.RunTasks = append(taskList.RunTasks, task.ParseRunTaskInfo())
		}
	}

	if reason == "startup" {
		tasks := model.GetTasksByStatus(model.TASK_STATUS_PENDING)
		for _, task := range tasks {
			taskList.RunTasks = append(taskList.RunTasks, task.ParseRunTaskInfo())
		}
	}

	logrus.Infof("client list task with reason:%s", reason)
	c.JSON(200, taskList)
}
func Running(c *gin.Context) {
	taskID := c.Query("taskId")
	task := model.GetTaskByUUID(taskID)
	start := c.Query("start")
	startedTs, err := strconv.Atoi(start)
	logrus.Infof("start:%v", startedTs)
	if err != nil {
		logrus.Errorf("parse start ts error: %v", err)
	}
	startedTm := time.Unix(int64(startedTs/1000), int64(startedTs%1000))

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
	if task.Status == model.TASK_STATUS_PENDING {
		task.StartedAt = &startedTm
		task.SetStatus(model.TASK_STATUS_RUNNING)
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

	start := c.Query("start")
	end := c.Query("end")
	exitCode := c.Query("exitCode")
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
	startedTs, _ := strconv.Atoi(start)
	endedTs, _ := strconv.Atoi(end)
	exitCodeint, _ := strconv.Atoi(exitCode)
	startedTm := time.Unix(int64(startedTs/1000), int64(startedTs%1000))
	endedTm := time.Unix(int64(startedTs/1000), int64(endedTs%1000))
	task.StartedAt = &startedTm
	task.EndedAt = &endedTm
	task.ExitCode = exitCodeint
	task.SetStatus(model.TASK_STATUS_FINISHED)
}
func Stopped(c *gin.Context) {
	taskID := c.Query("taskId")
	task := model.GetTaskByUUID(taskID)

	if task == nil {
		logrus.WithFields(logrus.Fields{
			"Module": "task",
			"Func":   "Stopped",
		}).Errorf("Get task failed by taskId:%s", taskID)
		c.AbortWithStatus(200)
		return
	}
	task.Status = model.TASK_STATUS_STOPPED
	model.GetDB().Save(task)
}
func Timeout(c *gin.Context) {
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
}
func Error(c *gin.Context) {
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
}
