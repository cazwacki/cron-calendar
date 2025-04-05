package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

func PutTask(context *gin.Context) {
	var task database.Task
	task.UserID = context.GetString("userId")
	task.ID = context.Param("taskId")
	if err := context.ShouldBindJSON(&task); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	if err := database.UpsertTask(task); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatus(http.StatusCreated)
}

func GetTask(context *gin.Context) {
	userId := context.GetString("userId")
	taskId := context.Param("taskId")
	task, err := database.GetTaskByUserAndId(userId, taskId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatusJSON(http.StatusOK, task)
}

func DeleteTask(context *gin.Context) {
	userId := context.GetString("userId")
	taskId := context.Param("taskId")
	if err := database.DeleteTaskByUserAndId(userId, taskId); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatus(http.StatusNoContent)
}

func GetAllTasks(context *gin.Context) {
	userId := context.GetString("userId")
	tasks, err := database.GetTasksByUserId(userId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatusJSON(http.StatusOK, tasks)
}
