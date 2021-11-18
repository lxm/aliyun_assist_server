package main

import (
	"github.com/lxm/aliyun_assist_server/pkg/apiserver"
	_ "github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func main() {
	model.ConnectDB()
	model.Migrate()
	// return
	router := apiserver.InitRouter()
	router.Run("0.0.0.0:10081")
	// router.RunTLS("0.0.0.0:443", "./aliyun-server.localdev02.qingjiao.link", "./aliyun-server.localdev02.qingjiao.link-key")
}
