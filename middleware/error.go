package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errortypes "cron-calendar/error_types"
)

func GenerateErrorResponse() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		error := context.Errors.Last().Err

		_, isAuthorizationError := error.(errortypes.AuthorizationError)
		if isAuthorizationError {
			context.AbortWithStatus(401)
			return
		}

		_, isValidationError := error.(errortypes.ValidationError)
		if isValidationError {
			context.AbortWithStatusJSON(400, gin.H{
				"error": error.Error(),
			})
			return
		}

		_, isDatabaseError := error.(errortypes.DatabaseError)
		if isDatabaseError {
			if error.Error() == "record not found" {
				context.AbortWithStatus(404)
				return
			}
			if strings.Contains(error.Error(), "constraint") {
				context.AbortWithStatus(409)
				return
			}
			context.AbortWithStatus(503)
			return
		}

		context.AbortWithStatus(http.StatusInternalServerError)
	}
}
