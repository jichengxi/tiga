package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"tiga-controller/ControllerTypes"
	"tiga-controller/cmd"
)

func main() {
	cmd.Execute()
	ControllerTypes.WorkGroup.Wait()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "以后再说")
	})
	r.GET("/daemonSet", func(c *gin.Context) {
		c.JSON(200, ControllerTypes.ResultData)
	})
	//log.Fatal(r.Run(":" + strconv.Itoa(ControllerTypes.Port)))
	log.Fatal(r.Run(":" + ControllerTypes.Port))

}
