package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userKey = "user_id"

// AuthRequired is a simple middleware to check the session.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

// GetUserID returns current user_id in session
func GetUserID(c *gin.Context) int64 {
	session := sessions.Default(c)
	user := session.Get(userKey)
	return user.(int64)
}

// SetUserID sets current user into session
func SetUserID(c *gin.Context, userID int64) error {
	session := sessions.Default(c)
	session.Set(userKey, userID)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

// ClearSession clears session info
func ClearSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		return errors.New("failed to save session")
	}
	return nil
}
