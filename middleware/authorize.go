package middleware

import (
	"crypto/sha3"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cron-calendar/database"
)

func ValidateAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionId := context.GetHeader("Session-ID")
		if sessionId == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		hash := sha3.New256()
		_, err := hash.Write([]byte(sessionId))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		finalHash := hash.Sum(nil)
		sessionId = string(finalHash)

		session, err := database.GetSessionById(sessionId)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if time.Now().After(session.ExpiresAt) {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("userId", session.UserID)
		context.Next()
	}
}
