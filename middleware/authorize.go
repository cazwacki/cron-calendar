package middleware

import (
	"crypto/sha3"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
)

var authorizeContext = log.With().Str("file", "authorize.go")

func ValidateAuthorization(db database.UserDB) gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := authorizeContext.Str("function", "ValidateAuthorization").Logger()
		functionLogger.Debug().Msg("invoked")

		sessionId := context.GetHeader("Session-ID")
		if sessionId == "" {
			functionLogger.Warn().Msg("attempt to hit auth-guarded endpoint without a session")
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		hash := sha3.New256()
		_, err := hash.Write([]byte(sessionId))
		if err != nil {
			functionLogger.Err(err).Msg("failed to write the generated session ID to hash")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		finalHash := hash.Sum(nil)
		sessionId = string(finalHash)

		session, err := db.GetSessionById(sessionId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to fetch user session from database")
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if session == nil {
			functionLogger.Debug().Msg("session not found")
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			functionLogger.Debug().Msg("session is expired")
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		functionLogger.Debug().Msg("returning with success")
		context.Set("userId", session.UserID)
		context.Next()
	}
}
