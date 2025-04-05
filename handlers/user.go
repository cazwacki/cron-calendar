package handlers

import (
	"crypto/rand"
	"crypto/sha3"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"cron-calendar/database"
	errortypes "cron-calendar/error_types"
)

func RegisterUser(context *gin.Context) {
	var user database.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	// salt and hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		context.Error(err)
		return
	}
	user.Password = string(hash)
	if err := database.InsertUser(user); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	context.AbortWithStatus(http.StatusCreated)
}

func Login(context *gin.Context) {
	var providedUser database.User
	if err := context.ShouldBindJSON(&providedUser); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	dbUser, err := database.GetUserById(providedUser.ID)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(providedUser.Password))
	if err != nil {
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
		context.Error(err)
		return
	}
	finalHash := hash.Sum(nil)
	session.ID = string(finalHash)

	if err = database.UpsertSession(session); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	// restore session ID for response
	session.ID = sessionId

	context.AbortWithStatusJSON(http.StatusOK, session)
}

func DestroyUser(context *gin.Context) {
	userId := context.GetString("userId")
	var providedUser database.User
	if err := context.ShouldBindJSON(&providedUser); err != nil {
		context.Error(errortypes.GenerateValidationError(err.Error()))
		return
	}
	if userId != providedUser.ID {
		context.Error(errortypes.GenerateAuthorizationError("attempt to delete other user"))
	}
	dbUser, err := database.GetUserById(providedUser.ID)
	if err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(providedUser.Password))
	if err != nil {
		context.Error(errortypes.GenerateAuthorizationError("invalid password"))
		return
	}

	if err := database.DeleteUserById(providedUser.ID); err != nil {
		context.Error(errortypes.GenerateDatabaseError(err.Error()))
		return
	}

	context.AbortWithStatus(http.StatusNoContent)
}
