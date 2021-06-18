package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/lxm/aliyun_assist_server/pkg/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	TASK_OUTPUT_STASH = "task:output:"
)

const (
	TASK_STATUS_PENDING  = "Pending"
	TASK_STATUS_RUNNING  = "Running"
	TASK_STATUS_FINISHED = "Finished"
	TASK_STATUS_FAIELD   = "Failed"
	TASK_STATUS_STOPPED  = "Stopped"
	TASK_STATUS_TIMEOUT  = "Timeout"
	TASK_STATUS_ERROR    = "Error"
)

type Task struct {
	ID         int            `json:"id" gorm:"primarykey"`
	UUID       string         `json:"task_id" gorm:"type:varchar(100);index"`
	BatchID    string         `json:"batch_id" grom:"varchar(100);index"`
	CommandID  int            `json:"command_id" gorm:"type:int"`
	InstanceID string         `json:"instance_id" gorm:"type:varchar(100);index"`
	Output     string         `json:"output" gorm:"type:text"`
	Status     string         `json:"string" gorm:"type:varchar(20);default:Pending"`
	TaskOption TaskOption     `gorm:"type:varchar(100);embedded"`
	ExitCode   int            `json:"exit_code" gorm:"type:tinyint"`
	StartedAt  *time.Time     `json:"started_at" grom:"default:NULL"`
	EndedAt    *time.Time     `json:"ended_at" grom:"default:NULL"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type TaskOption struct {
	RepeatMode string          `json:"repeat_mode" gorm:"type:varchar(100)"`
	Timed      bool            `json:"timed" gorm:"type:tinyint"`
	Frequency  string          `json:"frequency" gorm:"type:varchar(128)"`
	Parameters json.RawMessage `json:"parameters" gorm:"type:varchar(1024)"`
	Username   string          `json:"username" gorm:"type:varchar(100)"`
}

type RunTaskInfo struct {
	Task   Task       `json:"task"`
	Output OutputInfo `json:"output"`
	Repeat string     `json:"repeat"`
}

type SendFileTaskInfo struct {
	gorm.Model
	Content     string     `json:"content" gorm:"type:longtext"`
	ContentType string     `json:"contentType" gorm:"type:varchar(100)"`
	Destination string     `json:"destination" gorm:"type:varchar(1024)"`
	Group       string     `json:"group" gorm:"type:varchar(256);"`
	Mode        string     `json:"mode" gorm:"type:varchar(256);"`
	Name        string     `json:"name" gorm:"type:varchar(256);"`
	Overwrite   bool       `json:"overwrite" gorm:"type:int;"`
	Owner       string     `json:"owner" gorm:"type:varchar(256);"`
	Signature   string     `json:"signature" gorm:"type:varchar(256);"`
	TaskID      string     `json:"taskID" gorm:"type:varchar(256);index"`
	Timeout     int64      `json:"timeout" gorm:"type:int;"`
	Output      OutputInfo `json:"output" gorm:"column:output;type:longtext"`
}

type GshellCmd struct {
	Execute   string `json:"execute"`
	Arguments struct {
		Cmd string `json:"cmd"`
	} `json:"arguments"`
}

type GshellCmdReply struct {
	Return struct {
		CmdOutput string `json:"cmd_output"`
		Result    int    `json:"result"`
	} `json:"return"`
}

type OutputInfo struct {
	Interval  int  `json:"interval"`
	LogQuota  int  `json:"logQuota"`
	SkipEmpty bool `json:"skipEmpty"`
	SendStart bool `json:"sendStart"`
}

func (o *OutputInfo) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.New("decode outputinfo failed")
	}
	var output OutputInfo
	err := json.Unmarshal(bytes, &output)
	*o = output
	return err
}

func (o OutputInfo) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func GetTaskByID(taskID string) *Task {
	var task Task
	err := db.Where("id", taskID).Find(&task).Error
	if err != nil {
		return nil
	}
	return &task
}

func GetTaskByUUID(taskUUID string) *Task {
	var task Task
	err := db.Where("uuid", taskUUID).Find(&task).Error
	if err != nil {
		return nil
	}
	return &task
}

func CreateTask(commandID int, instanceId string, to TaskOption) *Task {
	var task Task
	task.CommandID = commandID
	task.InstanceID = instanceId
	task.TaskOption = to
	task.GenTaskUUID()

	err := db.Model(task).Save(&task).Error
	if err != nil {
		logrus.Errorf("CreateTask error %v", err)
		return nil
	}
	return &task
}

func (task *Task) GenTaskUUID() string {
	uuid := "t-" + util.RandStringRunes(32)
	task.UUID = uuid

	return uuid
}

func (task *Task) GetCommand() *Command {
	var command Command

	err := db.Where("id", task.CommandID).Find(&command).Error
	if err != nil {
		logrus.Errorf("Task GetCommand failed: %v", err)
		return nil
	}
	return &command
}

func (task *Task) ParseRunTaskInfo() map[string]interface{} {
	data := make(map[string]interface{})
	data["task"] = map[string]interface{}{
		"taskID":           task.UUID,
		"commandId":        task.GetCommand().UUID,
		"commandName":      task.GetCommand().Name,
		"commandContent":   task.GetCommand().CommandContent,
		"timeOut":          "60",
		"workingDirectory": "/tmp",
		"enableParameter":  false,
		"args":             "",
		"cron":             "",
		"type":             "RunShellScript",
	}
	data["output"] = map[string]interface{}{
		"sendStart": true,
		"skipEmpty": false,
		"logQuota":  102400,
		"interval":  1,
	}
	data["repeat"] = "Once"
	return data
}

func (task *Task) StashOutput(content string) error {
	stashKey := TASK_OUTPUT_STASH + task.UUID
	ctx := context.Background()
	redisClient := redisclient.GetClient()
	redisClient.RPush(ctx, stashKey, content)
	return nil
}

func (task *Task) DumpOutput() (string, error) {
	stashKey := TASK_OUTPUT_STASH + task.UUID
	ctx := context.Background()
	redisClient := redisclient.GetClient()
	outputLines, err := redisClient.LRange(ctx, stashKey, 0, -1).Result()
	if err != nil {
		return "", err
	}
	output := strings.Join(outputLines, "")
	task.Output = output
	err = db.Save(task).Error
	if err != nil {
		return "", err
	}
	return output, nil
}

func (task *Task) SetStatus(status string) error {
	task.Status = status
	return db.Save(task).Error
}
