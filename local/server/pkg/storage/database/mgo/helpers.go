package mgo

import (
	"fdim/pkg/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}
