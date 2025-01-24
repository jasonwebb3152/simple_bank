package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// Specially-named entry function for all tests in a go package
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run()) // Start running unit test
}
