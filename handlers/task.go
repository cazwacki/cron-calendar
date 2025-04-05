package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutTask(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	var task database.Task
	task.ID = context.Param("taskId")
	if err := context.ShouldBindJSON(&task); err != nil {
		context.Error(err)
		return
	}
	task.UserID = string(userId)
	if err := database.UpsertTask(task); err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusCreated)
}

func GetTask(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	taskId := context.Param("taskId")
	task, err := database.GetTaskByUserAndId(userId, taskId)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, task)
}

func DeleteTask(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	taskId := context.Param("taskId")
	if err := database.DeleteTaskByUserAndId(userId, taskId); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllTasks(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	tasks, err := database.GetTasksByUserId(userId)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, tasks)
}
