package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
)

var cronContext = log.With().Str("file", "cron.go")

func PutCronHandler(db database.CronDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := cronContext.Str("function", "PutCron").Logger()
		functionLogger.Debug().Msg("invoked")

		var cron database.Cron
		cron.UserID = context.GetString("userId")
		cron.ID = context.Param("cronId")
		if err := context.ShouldBindJSON(&cron); err != nil {
			functionLogger.Err(err).Msg("failed to bind payload to cron structure")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// todo: validate schedule
		if err := db.UpsertCronIfOwner(cron); err != nil {
			functionLogger.Err(err).Msg("failed to write cron to database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatus(http.StatusCreated)
	}
}

func GetCronHandler(db database.CronDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := cronContext.Str("function", "GetCron").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		cronId := context.Param("cronId")

		cron, err := db.GetCronByIdAndUserId(userId, cronId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to read cron from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if cron == nil {
			functionLogger.Debug().Msg("no cron found")
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatusJSON(http.StatusOK, &cron)
	}

}

func DeleteCronHandler(db database.CronDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := cronContext.Str("function", "DeleteCron").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		cronId := context.Param("cronId")
		if err := db.DeleteCronByIdAndUserId(userId, cronId); err != nil {
			functionLogger.Err(err).Msg("failed to delete cron from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatus(http.StatusNoContent)
	}

}

func GetAllCronsHandler(db database.CronDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := cronContext.Str("function", "GetAllCrons").Logger()
		functionLogger.Debug().Msg("invoked")

		userId := context.GetString("userId")
		crons, err := db.SearchCronsByUserId(userId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to get crons from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.AbortWithStatusJSON(http.StatusOK, crons)
	}
}
