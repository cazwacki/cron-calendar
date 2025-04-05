package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

func PutCron(context *gin.Context) {
	var cron database.Cron
	cron.UserID = context.GetString("userId")
	cron.ID = context.Param("cronId")
	if err := context.ShouldBindJSON(&cron); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	// todo: validate schedule
	if err := database.UpsertCron(cron); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatus(http.StatusCreated)
}

func GetCron(context *gin.Context) {
	userId := context.GetString("userId")
	cronId := context.Param("cronId")
	cron, err := database.GetCronByUserAndId(userId, cronId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatusJSON(http.StatusOK, cron)
}

func DeleteCron(context *gin.Context) {
	userId := context.GetString("userId")
	cronId := context.Param("cronId")
	if err := database.DeleteCronByUserAndId(userId, cronId); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatus(http.StatusNoContent)
}

func GetAllCrons(context *gin.Context) {
	userId := context.GetString("userId")
	crons, err := database.GetCronsByUserId(userId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatusJSON(http.StatusOK, crons)
}
