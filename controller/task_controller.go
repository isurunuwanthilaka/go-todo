package controller

import (
	"net/http"
	"strconv"
	"todo-app/utils/log"

	"todo-app/dto"
	"todo-app/service"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	log.Info("GetAllTasks request received")
	res, err := service.GetAllTasks()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	log.Info("GetTask request received", id)
	res, err := service.GetTask(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if res.Id == 0 {
		finalId, _ := strconv.Atoi(id)
		res.Id = uint(finalId)
		c.IndentedJSON(http.StatusNotFound, res)
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

func AddTask(c *gin.Context) {
	var task dto.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	log.Info("AddTask request received", task)
	res, err := service.AddTask(&task)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, res)
}
