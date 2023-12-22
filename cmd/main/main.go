package main

import (
	"assignment-2/app"

	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {

	router := gin.Default()

	v1 := router.Group("/v1/order")
	{
		v1.HEAD("/", app.Ping)
		v1.POST("/", app.CreateOrder)
		v1.GET("/:id", app.FetchOrder)
		v1.PUT("/:id", app.UpdateOrder)
	}

	return router
}

func main() {
	router := setUpRouter()
	router.Run()
}
