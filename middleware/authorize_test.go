package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"cron-calendar/middleware"
)

func TestValidateAuthorization_NoSessionId(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.ValidateAuthorization(), func(ctx *gin.Context) {
		t.Fail()
	})
	req := httptest.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()
	r.ServeHTTP(writer, req)

}
