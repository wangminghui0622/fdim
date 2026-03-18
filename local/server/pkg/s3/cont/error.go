package cont

import (
	"fmt"

	"fdim/pkg/s3"
)

type HashAlreadyExistsError struct {
	Object *s3.ObjectInfo
}

func (e *HashAlreadyExistsError) Error() string {
	return fmt.Sprintf("hash already exists: %s", e.Object.Key)
}
