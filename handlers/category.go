package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutCategory(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	var category database.Category
	category.ID = context.Param("categoryId")
	if err := context.ShouldBindJSON(&category); err != nil {
		context.Error(err)
		return
	}
	category.UserID = string(userId)
	if err := database.UpsertCategory(category); err != nil {
		context.Error(err)
		return
	}
	context.Status(http.StatusCreated)
}

func GetCategory(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	categoryId := context.Param("categoryId")
	category, err := database.GetCategoryByUserAndId(userId, categoryId)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusAccepted, category)
}

func DeleteCategory(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	categoryId := context.Param("categoryId")
	if err := database.DeleteCategoryByUserAndId(userId, categoryId); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllCategories(context *gin.Context) {
	value, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "cron-calendar: user not identified",
		})
		return
	}
	userId := fmt.Sprint(value)
	categories, err := database.GetCategoriesByUserId(userId)
	if err != nil {
		context.Error(err)
		return
	}
	context.JSON(http.StatusAccepted, categories)
}
