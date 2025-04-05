package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutCron(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	var cron database.Cron
	cron.ID = context.Param("cronId")
	if err := context.ShouldBindJSON(&cron); err != nil {
		context.Error(err)
		return
	}
	cron.UserID = string(userId)
	if err := database.UpsertCron(cron); err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusCreated)
}

func GetCron(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	cronId := context.Param("cronId")
	cron, err := database.GetCronByUserAndId(userId, cronId)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, cron)
}

func DeleteCron(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	cronId := context.Param("cronId")
	if err := database.DeleteCronByUserAndId(userId, cronId); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllCrons(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	crons, err := database.GetCronsByUserId(userId)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusOK, crons)
}
