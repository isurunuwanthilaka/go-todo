package app

import (
	"todo-app/controller"
	"todo-app/service"

	"github.com/gin-gonic/gin"
)

func StartApplication() {

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go service.WorkerRoutine()
	}

	router := gin.Default()

	router.GET("/tasks", controller.GetAllTasks)
	router.GET("/tasks/:id", controller.GetTask)
	router.POST("/tasks", controller.AddTask)

	router.Run(":8080")
}
