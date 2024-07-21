package authentication

import (
	"context"
	"net/http"
	"os"

	aws_conf "github.com/ank809/Chat-Application-golang/aws"
	"github.com/ank809/Chat-Application-golang/database"
	"github.com/ank809/Chat-Application-golang/helpers"
	"github.com/ank809/Chat-Application-golang/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Parse multipart form, set a reasonable memory limit
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Unable to parse form",
		})
		return
	}

	// Extract form values
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	phoneNumber := c.Request.FormValue("phoneNumber")
	about := c.Request.FormValue("about")

	// Check if required fields are provided
	if name == "" || email == "" || password == "" || phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Missing required fields",
		})
		return
	}

	//  profile picture
	profilePic, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Unable to upload file",
		})
		return
	}
	defer profilePic.Close()

	// Create user object
	user := models.User{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		About:       about,
	}

	// Validate password
	isValidPassword, res := helpers.CheckPassword(user.Password)
	if !isValidPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": res,
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)

	// Validate email
	isValidEmail, res := helpers.VerifyEmail(user.Email)
	if !isValidEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": res,
		})
		return
	}

	// Validate phone number
	isPhoneNumberValid, err := helpers.IsValidIndianPhoneNumber(user.PhoneNumber)
	if !isPhoneNumberValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Generate unique key for the profile picture
	uniqueKey := helpers.GetUniqueKey()
	imageUrl := user.Email + "/" + uniqueKey + header.Filename
	user.Profile_Url = imageUrl

	// Insert user into the database
	collection := database.OpenCollection(database.Client, "Users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to insert user",
		})
		return
	}

	// Upload profile picture to S3
	bucketName := os.Getenv("BUCKET_NAME")
	s3Client, err := aws_conf.GetS3Client()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to get S3 client",
		})
		return
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageUrl),
		Body:   profilePic,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to upload profile picture",
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account created successfully",
	})
}
