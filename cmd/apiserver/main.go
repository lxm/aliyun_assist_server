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
	router.RunTLS("0.0.0.0:443", "./aliyun-selfhosted.qingjiao.io.pem", "./aliyun-selfhosted.qingjiao.io-key.pem")
}
