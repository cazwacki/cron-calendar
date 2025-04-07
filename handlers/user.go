package handlers

import (
	"crypto/rand"
	"crypto/sha3"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

var userContext = log.With().Str("file", "user.go")

func RegisterUser(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "RegisterUser").Logger()
	functionLogger.Debug().Msg("invoked")

	var user database.User
	if err := context.ShouldBindJSON(&user); err != nil {
		functionLogger.Err(err).Msg("failed to bind payload to user structure")
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}

	// salt and hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		functionLogger.Err(err).Msg("failed to hash the password")
		context.Error(err)
		return
	}

	user.Password = string(hash)
	if err := database.InsertUser(user); err != nil {
		functionLogger.Err(err).Msg("failed to write user to database")
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusCreated)
}

func Login(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "Login").Logger()
	functionLogger.Debug().Msg("invoked")

	var providedUser database.User
	if err := context.ShouldBindJSON(&providedUser); err != nil {
		functionLogger.Err(err).Msg("failed to bind payload to user structure")
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}

	dbUser, err := database.GetUserById(providedUser.ID)
	if err != nil {
		functionLogger.Err(err).Msg("failed to get user from database")
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(providedUser.Password)); err != nil {
		functionLogger.Err(err).Msg("failed to compare the password to its hashed equivalent")
		context.Error(errortypes.GenerateAuthorizationError("invalid password"))
		return
	}

	// generate user session
	var session database.Session
	session.UserID = providedUser.ID
	session.ExpiresAt = time.Now().Add(time.Hour * 4)

	sessionId := rand.Text()
	hash := sha3.New256()
	_, err = hash.Write([]byte(sessionId))
	if err != nil {
		functionLogger.Err(err).Msg("failed to write the generated session ID to hash")
		context.Error(err)
		return
	}
	finalHash := hash.Sum(nil)
	session.ID = string(finalHash)

	if err := database.UpsertSession(session); err != nil {
		functionLogger.Err(err).Msg("failed to write new session to database")
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	// restore session ID for response
	session.ID = sessionId

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatusJSON(http.StatusOK, session)
}

func DestroyUser(context *gin.Context) {
	functionLogger := categoryContext.Str("function", "DestroyUser").Logger()
	functionLogger.Debug().Msg("invoked")

	userId := context.GetString("userId")
	var providedUser database.User
	if err := context.ShouldBindJSON(&providedUser); err != nil {
		functionLogger.Err(err).Msg("failed to bind payload to user structure")
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}

	if userId != providedUser.ID {
		functionLogger.Warn().Msgf("user %s attempted to delete user with id %s", userId, providedUser.ID)
		context.Error(errortypes.GenerateAuthorizationError("attempt to delete other user"))
		return
	}

	dbUser, err := database.GetUserById(providedUser.ID)
	if err != nil {
		functionLogger.Err(err).Msg("failed to get user from database")
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(providedUser.Password)); err != nil {
		functionLogger.Err(err).Msg("failed to compare the password to its hashed equivalent")
		context.Error(errortypes.GenerateAuthorizationError("invalid password"))
		return
	}

	if err := database.DeleteUserById(providedUser.ID); err != nil {
		functionLogger.Err(err).Msg("failed to delete user from database")
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	functionLogger.Debug().Msg("returning with success")
	context.AbortWithStatus(http.StatusNoContent)
}
