package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

var taskContext = log.With().Str("file", "task.go")

func PutTask(context *gin.Context) {
	functionLogger := taskContext.Str("function", "PutTask").Logger()
	functionLogger.Debug().Msg("invoked")

	var task database.Task
	task.UserID = context.GetString("userId")
	task.ID = context.Param("taskId")
	if err := context.ShouldBindJSON(&task); err != nil {
		functionLogger.Err(err).Msg("failed to bind payload to task structure")
		panic(errortypes.GenerateValidationError(err.Error()))
	}

	if err := database.UpsertTask(task); err != nil {
		functionLogger.Err(err).Msg("failed to write task to database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusCreated)
}

func GetTask(context *gin.Context) {
	functionLogger := taskContext.Str("function", "GetTask").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	taskId := context.Param("taskId")
	task, err := database.GetTaskByUserAndId(userId, taskId)
	if err != nil {
		functionLogger.Err(err).Msg("failed to read task from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatusJSON(http.StatusOK, task)
}

func DeleteTask(context *gin.Context) {
	functionLogger := taskContext.Str("function", "DeleteTask").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	taskId := context.Param("taskId")
	if err := database.DeleteTaskByUserAndId(userId, taskId); err != nil {
		functionLogger.Err(err).Msg("failed to delete task from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusNoContent)
}

func GetAllTasks(context *gin.Context) {
	functionLogger := taskContext.Str("function", "GetAllTasks").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	tasks, err := database.GetTasksByUserId(userId)
	if err != nil {
		functionLogger.Err(err).Msg("failed to get tasks from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatusJSON(http.StatusOK, tasks)
}
