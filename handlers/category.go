package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

var categoryContext = log.With().Str("file", "category.go")

func PutCategory(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "PutCategory").Logger()
	functionLogger.Debug().Msg("invoked")

	var category database.Category
	category.UserID = context.GetString("userId")
	category.ID = context.Param("categoryId")
	if err := context.ShouldBindJSON(&category); err != nil {
		functionLogger.Err(err).Msg("failed to bind payload to category structure")
		panic(errortypes.GenerateValidationError(err.Error()))
	}

	if err := database.UpsertCategory(category); err != nil {
		functionLogger.Err(err).Msg("failed to write category to database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusCreated)
}

func GetCategory(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "GetCategory").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	categoryId := context.Param("categoryId")
	category, err := database.GetCategoryByUserAndId(userId, categoryId)
	if err != nil {
		functionLogger.Err(err).Msg("failed to read category from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatusJSON(http.StatusAccepted, category)
}

func DeleteCategory(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "DeleteCategory").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	categoryId := context.Param("categoryId")
	if err := database.DeleteCategoryByUserAndId(userId, categoryId); err != nil {
		functionLogger.Err(err).Msg("failed to delete category from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusNoContent)
}

func GetAllCategories(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "GetAllCategories").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	categories, err := database.GetCategoriesByUserId(userId)
	if err != nil {
		functionLogger.Err(err).Msg("failed to get categories from database")
		panic(errortypes.GenerateDatabaseError(err.Error()))
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatusJSON(http.StatusAccepted, categories)
}
