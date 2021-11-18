package main

import (
	"github.com/lxm/aliyun_assist_server/pkg/manageserver"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func main() {
	model.ConnectDB()
	model.Migrate()
	router := manageserver.InitRouter()
	router.Run("0.0.0.0:18080")
}
