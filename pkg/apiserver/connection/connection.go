package connection

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/types"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
	"github.com/sirupsen/logrus"
)

func Detect(c *gin.Context) {}

func Gshell(c *gin.Context) {}

func HeartBeat(c *gin.Context) {
	resp := types.HeartBeatResp{
		NextInterval: 10,
		NewTasks:     false,
	}
	c.JSON(200, resp)
}

func NotifyServer(c *gin.Context) {
	instanceID := c.GetString("checked-instance-id")
	logrus.Infof("NotifyServer-Start serve ws for instance: %v", instanceID)
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
	go func() {
		ctx := context.Background()
		pubsub := redisClient.Subscribe(ctx, "notify_server:"+instanceID)
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
	logrus.Infof("NotifyServer-End serve ws for instance: %v", instanceID)
}

func PluginList(c *gin.Context) {}

func ExceptionClientReport(c *gin.Context) {
	rawData, _ := c.GetRawData()
	logrus.Infof("ExceptionClientReport:%v", string(rawData))
}
