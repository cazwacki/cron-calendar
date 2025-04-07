package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	errortypes "cron-calendar/error_types"
)

var errorContext = log.With().Str("file", "error.go")

func GenerateErrorResponse() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		functionLogger := errorContext.Str("function", "GenerateErrorResponse").Logger()
		functionLogger.Debug().Msg("invoked")

		error := context.Errors.Last().Err
		functionLogger.Debug().Msgf("given error: %s", error.Error())

		_, isAuthorizationError := error.(errortypes.AuthorizationError)
		if isAuthorizationError {
			functionLogger.Debug().Msg("returning authorization error")
			context.AbortWithStatus(401)
			return
		}

		_, isValidationError := error.(errortypes.ValidationError)
		if isValidationError {
			functionLogger.Debug().Msg("returning validation error")
			context.AbortWithStatusJSON(400, gin.H{
				"error": error.Error(),
			})
			return
		}

		_, isDatabaseError := error.(errortypes.DatabaseError)
		if isDatabaseError {
			if error.Error() == "record not found" {
				functionLogger.Debug().Msg("returning not found database error")
				context.AbortWithStatus(404)
				return
			}
			if strings.Contains(error.Error(), "constraint") {
				functionLogger.Debug().Msg("returning constraint database error")
				context.AbortWithStatus(409)
				return
			}
			functionLogger.Debug().Msg("returning generic database error")
			context.AbortWithStatus(503)
			return
		}

		functionLogger.Debug().Msg("returning internal server error")
		context.AbortWithStatus(http.StatusInternalServerError)
	}
}
