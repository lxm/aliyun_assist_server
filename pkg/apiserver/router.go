package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/connection"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/instance"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/metrics"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/session"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/task"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver/update"
	_ "github.com/lxm/aliyun_assist_server/pkg/config"
)

const (
	URI_PREFIX = "/api"
)

/*
/luban/api/v1/update/update_check
/luban/api/connection_detect
/luban/api/v1/task/invalid
/luban/api/v1/task/list
/luban/api/v1/task/running
/luban/api/v1/task/finish
/luban/api/v1/task/stopped
/luban/api/v1/task/timeout
/luban/api/v1/task/error
/luban/api/heart-beat
/luban/api/gshell
/luban/api/v1/plugin/list
/luban/api/v1/exception/client_report
/luban/api/instance/register
/luban/api/instance/deregister
*/
func InitRouter(prefix string, r *gin.Engine) *gin.Engine {

	r.Use(gin.Logger())
	r.RemoveExtraSlash = true
	lubanGroup := r.Group(prefix+"/luban", instance.CheckHeaderMiddleware)

	lubanGroup.GET("/notify_server", connection.NotifyServer)
	lubanGroup.POST(URI_PREFIX+"/v1/update/update_check", update.Check)
	lubanGroup.POST(URI_PREFIX+"/v1/exception/client_report", connection.ExceptionClientReport)
	lubanGroup.GET(URI_PREFIX+"/heart-beat", connection.HeartBeat)
	lubanGroup.POST(URI_PREFIX+"/instance/register", instance.Reg)
	lubanGroup.POST(URI_PREFIX+"/metrics", metrics.Metrics)

	apiGroupTask := lubanGroup.Group(URI_PREFIX + "/v1/task")
	apiGroupTask.POST("/list", task.List)
	apiGroupTask.POST("/running", task.Running)
	apiGroupTask.POST("/finish", task.Finish)
	apiGroupTask.POST("/stopped", task.Stopped)
	apiGroupTask.POST("/timeout", task.Timeout)
	apiGroupTask.POST("/error", task.Error)
	apiGroupTask.POST("/invalid", task.Invalid)

	sessionGroup := lubanGroup.Group(URI_PREFIX + "/v1/session")

	sessionGroup.GET("/list", session.List)
	sessionGroup.Any("/websocket", session.Websocket)

	return r

}
