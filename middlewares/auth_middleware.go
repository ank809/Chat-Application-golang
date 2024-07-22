package middlewares

import (
	"net/http"
	"os"
	"time"

	"github.com/ank809/Chat-Application-golang/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing cookie"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}

		claims := &models.Claims{}
		JWT_KEY := []byte(os.Getenv("JWT_KEY"))
		token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			return JWT_KEY, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Refresh token if it is about to expire
		expirationTime := time.Unix(claims.ExpiresAt, 0)
		if time.Until(expirationTime) < 5*time.Second {
			newExpirationTime := time.Now().Add(10 * time.Minute)
			claims.ExpiresAt = newExpirationTime.Unix()
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := newToken.SignedString(JWT_KEY)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			http.SetCookie(c.Writer, &http.Cookie{
				Name:    "token",
				Value:   tokenStr,
				Expires: newExpirationTime,
			})
		}
		c.Set("user", claims)
		c.Next()
	}
}
