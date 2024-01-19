package metrics

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Report(c *gin.Context) {
	jsonData, _ := io.ReadAll(c.Request.Body)
	logrus.Infof("metrics:" + string(jsonData))
}
