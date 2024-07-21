package helpers

import (
	"context"

	emailVerifier "github.com/AfterShip/email-verifier"
	"github.com/ank809/Chat-Application-golang/database"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	verifier = emailVerifier.NewVerifier()
)

func VerifyEmail(email string) (bool, string) {
	res, err := verifier.Verify(email)
	if err != nil {
		return false, "verify email address failed"
	}
	if !res.Syntax.Valid {
		return false, "Invalid email syntax"
	}

	coll := database.OpenCollection(database.Client, "Users")
	coll.FindOne(context.TODO(), bson.M{"email": email}).Err()
	err = coll.FindOne(context.TODO(), bson.M{"email": email}).Err()
	if err == nil {
		return false, "Email already in use"
	} else if err != mongo.ErrNoDocuments {
		return false, "Error checking email existence"
	}

	return true, "Valid email"
}
