package middleware

import (
	"crypto/sha3"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

var authorizeContext = log.With().Str("file", "authorize.go")

func ValidateAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		functionLogger := authorizeContext.Str("function", "ValidateAuthorization").Logger()
		functionLogger.Debug().Msg("invoked")

		sessionId := context.GetHeader("Session-ID")
		if sessionId == "" {
			functionLogger.Warn().Msg("attempt to hit auth-guarded endpoint without a session")
			panic(errortypes.GenerateAuthorizationError("attempt to hit auth-guarded endpoint without a session"))
		}
		hash := sha3.New256()
		_, err := hash.Write([]byte(sessionId))
		if err != nil {
			functionLogger.Err(err).Msg("failed to write the generated session ID to hash")
			panic(errors.New("failed to write the generated session ID to hash"))
		}
		finalHash := hash.Sum(nil)
		sessionId = string(finalHash)

		session, err := database.GetSessionById(sessionId)
		if err != nil {
			functionLogger.Err(err).Msg("failed to fetch user session from database")
			panic(errortypes.GenerateDatabaseError(err.Error()))
		}
		if time.Now().After(session.ExpiresAt) {
			functionLogger.Debug().Msg("session is expired")
			panic(errortypes.GenerateAuthorizationError("session is expired"))
		}

		functionLogger.Debug().Msg("returning with success")
		context.Set("userId", session.UserID)
		context.Next()
	}
}
