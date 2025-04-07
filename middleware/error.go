package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	errortypes "cron-calendar/error_types"
)

var errorContext = log.With().Str("file", "error.go")

func GenerateErrorResponse(context *gin.Context, received any) {
	functionLogger := errorContext.Str("function", "GenerateErrorResponse").Logger()
	functionLogger.Debug().Msg("invoked")

	err, ok := received.(error)

	if !ok {
		functionLogger.Warn().Msg("recovered from panic without an error provided")
		context.AbortWithStatus(500)
	}

	functionLogger.Debug().Msgf("given error: %s", err.Error())

	_, isAuthorizationError := err.(errortypes.AuthorizationError)
	if isAuthorizationError {
		functionLogger.Debug().Msg("returning authorization error")
		context.AbortWithStatus(401)
		return
	}

	_, isValidationError := err.(errortypes.ValidationError)
	if isValidationError {
		functionLogger.Debug().Msg("returning validation error")
		context.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, isDatabaseError := err.(errortypes.DatabaseError)
	if isDatabaseError {
		if err.Error() == "record not found" {
			functionLogger.Debug().Msg("returning not found database error")
			context.AbortWithStatus(404)
			return
		}
		if strings.Contains(err.Error(), "constraint") {
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
