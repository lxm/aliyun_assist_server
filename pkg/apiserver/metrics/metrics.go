package metrics

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Metrics(c *gin.Context) {
	jsonData, _ := io.ReadAll(c.Request.Body)
	logrus.Infof("metrics:%s", string(jsonData))
	c.JSON(200, gin.H{})
}
