package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/manageserver"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func main() {
	config.LoadConfig()
	model.ConnectDB()
	model.Migrate()
	r := gin.New()
	manageserver.InitRouter("/", r)
	r.Run("0.0.0.0:18080")
}
