package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

func PutCategory(context *gin.Context) {
	var category database.Category
	category.UserID = context.GetString("userId")
	category.ID = context.Param("categoryId")
	if err := context.ShouldBindJSON(&category); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	if err := database.UpsertCategory(category); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatus(http.StatusCreated)
}

func GetCategory(context *gin.Context) {
	userId := context.GetString("userId")
	categoryId := context.Param("categoryId")
	category, err := database.GetCategoryByUserAndId(userId, categoryId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatusJSON(http.StatusAccepted, category)
}

func DeleteCategory(context *gin.Context) {
	userId := context.GetString("userId")
	categoryId := context.Param("categoryId")
	if err := database.DeleteCategoryByUserAndId(userId, categoryId); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatus(http.StatusNoContent)
}

func GetAllCategories(context *gin.Context) {
	userId := context.GetString("userId")
	categories, err := database.GetCategoriesByUserId(userId)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatusJSON(http.StatusAccepted, categories)
}
