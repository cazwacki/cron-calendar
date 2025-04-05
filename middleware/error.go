package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GenerateErrorResponse() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		error := context.Errors[0]

		if error.Error() == "record not found" {
			context.Status(404)
			return
		}

		if strings.Contains(error.Error(), "Field validation") {
			context.JSON(400, error.JSON())
			return
		}

		context.JSON(http.StatusInternalServerError, error.JSON())
	}
}
