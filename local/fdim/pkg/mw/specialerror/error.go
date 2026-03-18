package specialerror

import (
	"errors"

	"fdim/pkg/errs"
)

var handlers []func(err error) errs.CodeError

func AddErrHandler(h func(err error) errs.CodeError) (err error) {
	if h == nil {
		return errs.New("nil handler")
	}
	handlers = append(handlers, h)
	return nil
}

func AddReplace(target error, codeErr errs.CodeError) error {
	handler := func(err error) errs.CodeError {
		if errors.Is(err, target) {
			return codeErr
		}
		return nil
	}

	if err := AddErrHandler(handler); err != nil {
		return err
	}

	return nil
}

func ErrCode(err error) errs.CodeError {
	var codeErr errs.CodeError
	if errors.As(err, &codeErr) {
		return codeErr
	}
	for i := 0; i < len(handlers); i++ {
		if codeErr := handlers[i](err); codeErr != nil {
			return codeErr
		}
	}
	return nil
}
