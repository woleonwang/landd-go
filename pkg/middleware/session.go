package middleware

import (
	"errors"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userKey = "user"

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

func AdminAuth(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if user.(mysql.User).Role != model.Admin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized as admin"})
		return
	}
	c.Next()
}

func RecruiterAuth(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if user.(mysql.User).Role == model.Admin {
		c.Next()
		return
	}
	if user.(mysql.User).Role != model.Recruiter {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized as recruiter"})
		return
	}
	c.Next()
}

func PartnerAuth(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if user.(mysql.User).Role == model.Admin {
		c.Next()
		return
	}
	if user.(mysql.User).Role != model.Partner {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized as partner"})
		return
	}
	c.Next()
}

// GetUser returns current user in session
func GetUser(c *gin.Context) mysql.User {
	session := sessions.Default(c)
	user := session.Get(userKey)
	return user.(mysql.User)
}

// SetUser sets current user into session
func SetUser(c *gin.Context, user mysql.User) error {
	session := sessions.Default(c)
	session.Set(userKey, user)
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
