package dbutil

import (
	"fdim/pkg/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsDBNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}
