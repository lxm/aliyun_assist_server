package manageserver

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/manageserver/activation"
	"github.com/lxm/aliyun_assist_server/pkg/manageserver/command"
	"github.com/lxm/aliyun_assist_server/pkg/manageserver/invocation"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func InitRouter(prefix string, r *gin.Engine, middleware ...gin.HandlerFunc) *gin.Engine {
	model.ConnectDB()
	r.Use(middleware...)
	r.RemoveExtraSlash = true
	manageGroup := r.Group(prefix)

	//TODO command manage
	// manageGroup.POST("/command")
	// manageGroup.DELETE("/command/:id")
	// manageGroup.GET("/command/:id")
	// manageGroup.PUT("/command/:id")

	manageGroup.POST("/activationcode", activation.CreateActivationCode)
	manageGroup.GET("/activationcode/:code/instances", activation.GetInstancesByActivationCode)
	manageGroup.GET("/activationcode/instances", activation.BatchGetInstancesByActivationCode)

	manageGroup.POST("/command/:id/invoke", command.InvokeCommand)
	manageGroup.POST("/command/run", command.RunCommand)

	// manageGroup.GET("/invocations", invocation.ListInvocations)
	manageGroup.GET("/invocationresults", invocation.ListInvocationResults)

	//TODO sendfile
	// manageGroup.POST("/sendfile")
	// manageGroup.GET("/instance")

	return r
}
