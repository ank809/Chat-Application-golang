package authentication

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func Login(c *gin.Context) {

	var user models.LoginUser
	var foundUser models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email and password are required",
		})
		return
	}

	// Find in database

	coll := database.OpenCollection(database.Client, "Users")
	err = coll.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid email ",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid password",
		})
		return
	}

	token, err := generateJWT(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate token",
		})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Minute * 10),
	})
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
}

func generateJWT(user models.User) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	claims := models.Claims{
		Email: user.Email,
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))

}
