package main

import (
	"context"
	"fmt"

	_ "github.com/lxm/aliyun_assist_server/pkg/config"
	"github.com/lxm/aliyun_assist_server/pkg/redisclient"
)

func main() {
	redisClient := redisclient.GetClient()
	instanceID := "i-p8YjI2Uk3X6SCZe0"
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, "notify_server:"+instanceID)
	for msg := range pubsub.Channel() {
		fmt.Println(msg.Channel, msg.Payload, "\r\n")
	}
}
