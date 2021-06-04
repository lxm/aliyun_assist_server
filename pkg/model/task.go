package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Task struct {
	ID         int            `json:"id" gorm:"primarykey"`
	CommandID  int            `json:"command_id" gorm:"type:int"`
	TaskOption TaskOption     `gorm:"type:varchar(100);embedded"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type TaskOption struct {
	InstanceIDs string          `json:"instance_id"`
	RepeatMode  string          `json:"repeat_mode" gorm:"type:varchar(100)"`
	Timed       bool            `json:"timed" gorm:"type:tinyint"`
	Frequency   string          `json:"frequency" gorm:"type:varchar(128)"`
	Parameters  json.RawMessage `json:"parameters" gorm:"type:varchar(1024)"`
	Username    string          `json:"username" gorm:"type:varchar(100)"`
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
	Destination string     `json:"destination" gorm:"type:varchar(1024);`
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

func CreateTask(commandID int, to TaskOption) *Task {
	var task Task
	task.CommandID = commandID
	task.TaskOption = to

	err := db.Model(task).Save(&task).Error
	if err != nil {
		logrus.Errorf("CreateTask error %v", err)
		return nil
	}
	return &task
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
		"taskID":           task.ID,
		"commandId":        task.CommandID,
		"commandContent":   task.GetCommand().CommandContent,
		"timeOut":          60,
		"workingDirectory": "/tmp",
		"enableParameter":  false,
		"args":             "",
		"cron":             "",
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
