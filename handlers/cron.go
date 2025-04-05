package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutCron(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	var cron database.Cron
	cron.ID = context.Param("cronId")
	if err := context.ShouldBindJSON(&cron); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cron.UserID = string(userId)
	if err := database.UpsertCron(cron); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.Status(http.StatusCreated)
}

func GetCron(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	cronId := context.Param("cronId")
	cron, err := database.GetCronByUserAndId(userId, cronId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, cron)
}

func DeleteCron(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	cronId := context.Param("cronId")
	if err := database.DeleteCronByUserAndId(userId, cronId); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllCrons(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	crons, err := database.GetCronsByUserId(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, crons)
}
