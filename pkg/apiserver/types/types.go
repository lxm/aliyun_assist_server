package types

type HeartBeatResp struct {
	NextInterval float64 `json:"nextInterval"`
	NewTasks     bool    `json:"newTasks"`
}

type TaskListResp struct {
	RunTasks      []interface{} `json:"run"`
	StopTasks     []interface{} `json:"stop"`
	SendFileTasks []interface{} `json:"file"`
	SessionTasks  []interface{} `json:"session"`
	InstanceId    string        `json:"instanceId"`
}
