package invocation

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lxm/aliyun_assist_server/pkg/model"
)

func ListInvocations(c *gin.Context) {

}

/*
当前仅实现根据invoke_id查询
可以支持多个，使用，隔开
*/
func ListInvocationResults(c *gin.Context) {

	invokeIDRaw := c.Query("invoke_id")
	invokeIDs := strings.Split(invokeIDRaw, ",")
	if len(invokeIDs) == 0 {
		c.JSON(200, gin.H{})
		return
	}

	tasks, err := model.ListTasksByInvokeIDs(invokeIDs)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"results": tasks,
	})

}
