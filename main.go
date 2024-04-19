package main

import (
	"github.com/gin-gonic/gin"
	"go-crud/controllers"
	"go-crud/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}
func main() {
	r := gin.Default()
	r.POST("/post", controllers.PostCreate)
	r.GET("/post", controllers.PostIndex)
	r.GET("/post/:id", controllers.PostShow)
	r.PUT("/post/:id", controllers.PostUpdate)
	r.DELETE("/post/:id", controllers.PostDelete)

	r.Run() // listen and serve on 0.0.0.0:8080
}
