package connection

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/types"
	"github.com/lxm/aliyun_assist_server/pkg/model"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/sirupsen/logrus"
)

func Detect(c *gin.Context) {}

func Gshell(c *gin.Context) {}

func HeartBeat(c *gin.Context) {
	instanceID := c.GetString("checked-instance-id")
	// logrus.Infof("HeartBeat:%v", instanceID)
	// get pengding tasks
	pendingTasks, err := model.ListInstancePendingTasks(instanceID)
	newTasks := false
	var nextInterval float64 = 10
	if err != nil {
		logrus.Errorf("HeartBeat get pending tasks failed: %v", err)
	} else {
		if len(pendingTasks) > 0 {
			redisClient := redisclient.GetClient()
			for _, task := range pendingTasks {
				num, err := task.SendKickMsg(redisClient)
				logrus.Infof("send kick msg to %v, num:%v, err: %v", task.InstanceID, num, err)
			}
			nextInterval = (float64)(10 * len(pendingTasks))
			newTasks = true
		}
	}
	resp := types.HeartBeatResp{
		NextInterval: nextInterval,
		NewTasks:     newTasks,
	}
	c.JSON(200, resp)
}

func NotifyServer(c *gin.Context) {
	instanceID := c.GetString("checked-instance-id")
	logrus.Infof("NotifyServer-Start serve ws for instance: %v", instanceID)
	instance := model.GetIntanceByID(instanceID)
	if instance == nil {
		logrus.Errorf("NotifyServer instance not found with instance id: [%s]", instanceID)
	}
	conn := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := conn.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("NotifyServer Upgrade err:%v", err)
	}
	defer ws.Close()
	wsQuit := make(chan int)
	redisClient := redisclient.GetClient()
	var writeMu sync.Mutex
	//Process msg
	instance.SetOnline()
	go func() {
		ctx := context.Background()
		pubsub := redisClient.Subscribe(ctx, "notify_server:"+instanceID)
		// send <- struct{}{}
		defer pubsub.Close()
	msgProcess:
		for msg := range pubsub.Channel() {
			logrus.Infof("receive msg:%v", msg)
			select {
			case <-wsQuit:
				break msgProcess
			default:
				writeMu.Lock()
				ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				writeMu.Unlock()
				logrus.Infof("msg:%v", msg)
			}
		}
	}()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		logrus.Infof("NotifyServer receive mt:%v msg:%v", mt, string(message))
	}
	if err != nil {
		logrus.Errorf("NotifyServer serve ws failed: %v", err)
	}
	wsQuit <- 1
	instance.SetOffline()
	logrus.Infof("NotifyServer-End serve ws for instance: %v", instanceID)
}

func PluginList(c *gin.Context) {}

func ExceptionClientReport(c *gin.Context) {
	rawData, _ := c.GetRawData()
	logrus.Infof("ExceptionClientReport:%v", string(rawData))
}
