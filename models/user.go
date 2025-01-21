package models 

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Email string			 `bson:"email" json:"email"`
	Password string 		 `bson:"password" json:"password"`
	isVerified bool		 `bson: "isVerified" json:"isVerified"`
}



