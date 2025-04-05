package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutTask(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	var task database.Task
	task.ID = context.Param("taskId")
	if err := context.ShouldBindJSON(&task); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task.UserID = string(userId)
	if err := database.UpsertTask(task); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.Status(http.StatusCreated)
}

func GetTask(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	taskId := context.Param("taskId")
	task, err := database.GetTaskByUserAndId(userId, taskId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, task)
}

func DeleteTask(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	taskId := context.Param("taskId")
	if err := database.DeleteTaskByUserAndId(userId, taskId); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllTasks(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	tasks, err := database.GetTasksByUserId(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, tasks)
}
