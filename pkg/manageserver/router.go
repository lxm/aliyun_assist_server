package manageserver

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/manageserver/command"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	model.ConnectDB()
	manageGroup := r.Group("/")

	// manageGroup.POST("/command")
	// manageGroup.DELETE("/command/:id")
	// manageGroup.GET("/command/:id")
	// manageGroup.PUT("/command/:id")

	// manageGroup.POST("/command/:id/invoke")

	manageGroup.POST("/command/run", command.RunCommand)

	// manageGroup.POST("/sendfile")

	// manageGroup.GET("/instance")

	return r
}
