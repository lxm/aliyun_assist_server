package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/apiserver"
	"github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func main() {
	config.LoadConfig()
	model.ConnectDB()
	model.Migrate()
	// return
	r := gin.New()
	apiserver.InitRouter("/", r)
	go func() {
		r.RunTLS("0.0.0.0:10083", "./docker.for.mac.localhost.pem", "./docker.for.mac.localhost-key.pem")

	}()
	r.Run("0.0.0.0:10081")
}
