package controller

import "github.com/gin-gonic/gin"

func DemoController(router *gin.Engine) {
	router.GET("/ping", ping)
	router.GET("/", hello)
	router.POST("/ping", pingpost)
	router.GET("/:name", pingname)

}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"massage": "pong!!!!",
	})
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"massage": "Hello Mikkee",
	})
}
func pingpost(c *gin.Context) {
	name := c.PostForm("name")
	nickname := c.DefaultPostForm("nickname", "Nodata")
	c.JSON(200, gin.H{
		"Name":     "post " + name,
		"Nickname": "post " + nickname,
	})
}
func pingname(c *gin.Context) {
	username := c.Param("name")
	age := c.Query("age")
	c.JSON(200, gin.H{
		"UserName": username + " " + age,
	})
}
