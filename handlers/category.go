package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
)

var categoryContext = log.With().Str("file", "category.go")

func PutCategoryHandler(db database.CategoryDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := categoryContext.Str("function", "PutCategory").Logger()
		functionLogger.Debug().Msg("invoked")

		var category database.Category
		category.UserID = context.GetString("userId")
		category.ID = context.Param("categoryId")
		if err := context.ShouldBindJSON(&category); err != nil {
			functionLogger.Err(err).Msg("failed to bind payload to category structure")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.UpsertCategoryIfOwner(category); err != nil {
			functionLogger.Err(err).Msg("failed to write category to database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatus(http.StatusCreated)
	}
}

func GetCategoryHandler(db database.CategoryDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := categoryContext.Str("function", "GetCategory").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		categoryId := context.Param("categoryId")

		category, err := db.GetCategoryByIdAndUserId(categoryId, userId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to read category from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if category == nil {
			functionLogger.Debug().Msg("no category found")
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatusJSON(http.StatusAccepted, &category)
	}
}

func DeleteCategoryHandler(db database.CategoryDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := categoryContext.Str("function", "DeleteCategory").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		categoryId := context.Param("categoryId")
		if err := db.DeleteCategoryByIdAndUserId(categoryId, userId); err != nil {
			functionLogger.Err(err).Msg("failed to delete category from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatus(http.StatusNoContent)
	}
}

func GetAllCategoriesHandler(db database.CategoryDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := categoryContext.Str("function", "GetAllCategories").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		categories, err := db.SearchCategoriesByUserId(userId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to get categories from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatusJSON(http.StatusAccepted, categories)
	}
}
