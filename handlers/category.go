package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func PutCategory(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	var category database.Category
	category.ID = context.Param("categoryId")
	if err := context.ShouldBindJSON(&category); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	category.UserID = string(userId)
	if err := database.UpsertCategory(category); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.Status(http.StatusCreated)
}

func GetCategory(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	categoryId := context.Param("categoryId")
	category, err := database.GetCategoryByUserAndId(userId, categoryId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusAccepted, category)
}

func DeleteCategory(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	categoryId := context.Param("categoryId")
	if err := database.DeleteCategoryByUserAndId(userId, categoryId); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.Status(http.StatusNoContent)
}

func GetAllCategories(context *gin.Context) {
	userId := fmt.Sprint(context.Get("userId"))
	categories, err := database.GetCategoriesByUserId(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusAccepted, categories)
}
