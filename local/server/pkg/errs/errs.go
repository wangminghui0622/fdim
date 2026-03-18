package errs

import (
	"bytes"
	"errors"
	"fmt"
)

// UnknownCode represents the error code when code is not parsed or parsed code equals 0.
const UnknownCode = 1000

// Chat/Account error codes.
const (
	FormattingError      = 10001 // Error in formatting
	HasRegistered        = 10002 // user has already registered
	NotRegistered        = 10003 // user is not registered
	PasswordErr          = 10004 // Password error
	GetIMTokenErr        = 10005 // Error in getting IM token
	RepeatSendCode       = 10006 // Repeat sending code
	MailSendCodeErr      = 10007 // Error in sending code via email
	SmsSendCodeErr       = 10008 // Error in sending code via SMS
	CodeInvalidOrExpired = 10009 // Code is invalid or expired
	RegisterFailed       = 10010 // Registration failed
	ResetPasswordFailed  = 10011 // Resetting password failed
	RegisterLimit        = 10012 // Registration limit exceeded
	LoginLimit           = 10013 // Login limit exceeded
	InvitationError      = 10014 // Error in invitation
)

// General error codes.
const (
	NoError = 0 // No error

	DatabaseError = 90002 // Database error (redis/mysql, etc.)
	NetworkError  = 90004 // Network error
	DataError     = 90007 // Data error

	CallbackError = 80000

	// General error codes.
	ServerInternalError   = 500  // Server internal error
	ArgsError             = 1001 // Input parameter error
	NoPermissionError     = 1002 // Insufficient permission
	DuplicateKeyError     = 1003
	RecordNotFoundError   = 1004 // Record does not exist
	SecretNotChangedError = 1050 // secret not changed

	// Account error codes.
	UserIDNotFoundError    = 1101 // UserID does not exist or is not registered
	RegisteredAlreadyError = 1102 // user is already registered

	// Group error codes.
	GroupIDNotFoundError  = 1201 // GroupID does not exist
	GroupIDExisted        = 1202 // GroupID already exists
	NotInGroupYetError    = 1203 // Not in the group yet
	DismissedAlreadyError = 1204 // Group has already been dismissed
	GroupTypeNotSupport   = 1205
	GroupRequestHandled   = 1206

	// Relationship error codes.
	CanNotAddYourselfError   = 1301 // Cannot add yourself as a friend
	BlockedByPeer            = 1302 // Blocked by the peer
	NotPeersFriend           = 1303 // Not the peer's friend
	RelationshipAlreadyError = 1304 // Already in a friend relationship

	// Message error codes.
	MessageHasReadDisable = 1401
	MutedInGroup          = 1402 // Member muted in the group
	MutedGroup            = 1403 // Group is muted
	MsgAlreadyRevoke      = 1404 // Message already revoked
	LicenseExpired        = 1405 // License expired

	// Token error codes.
	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507

	// Long connection gateway error codes.
	ConnOverMaxNumLimit  = 1601
	ConnArgsErr          = 1602
	PushMsgErr           = 1603
	IOSBackgroundPushErr = 1604

	// S3 error codes.
	FileUploadedExpiredError = 1701 // Upload expired
)

// CodeError 错误代码接口
type CodeError interface {
	error
	Code() int
	Msg() string
	Detail() string
	Wrap() error
	WrapMsg(msg string, kv ...any) error
	WithDetail(detail string) CodeError
	Is(err error) bool
}

// codeError 错误代码实现
type codeError struct {
	code   int
	msg    string
	detail string
}

func (e *codeError) Error() string {
	if e.detail != "" {
		return fmt.Sprintf("%d %s %s", e.code, e.msg, e.detail)
	}
	return fmt.Sprintf("%d %s", e.code, e.msg)
}

func (e *codeError) Code() int {
	return e.code
}

func (e *codeError) Msg() string {
	return e.msg
}

func (e *codeError) Detail() string {
	return e.detail
}

func (e *codeError) Wrap() error {
	return e
}

func (e *codeError) WrapMsg(msg string, kv ...any) error {
	newErr := &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: e.detail,
	}
	if msg != "" || len(kv) > 0 {
		detail := toString(msg, kv)
		if newErr.detail == "" {
			newErr.detail = detail
		} else {
			newErr.detail += ", " + detail
		}
	}
	return newErr
}

func (e *codeError) WithDetail(detail string) CodeError {
	var d string
	if e.detail == "" {
		d = detail
	} else {
		d = e.detail + ", " + detail
	}
	return &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: d,
	}
}

func (e *codeError) Is(err error) bool {
	if err == nil {
		return e == nil
	}
	var codeErr CodeError
	if errors.As(Unwrap(err), &codeErr) {
		return e.code == codeErr.Code()
	}
	return false
}

// NewCodeError 创建新的错误代码
func NewCodeError(code int, msg string) CodeError {
	return &codeError{
		code: code,
		msg:  msg,
	}
}

// New 创建新的错误
func New(msg string, kv ...any) CodeError {
	return &codeError{
		code:   ServerInternalError,
		msg:    toString(msg, kv),
		detail: "",
	}
}

// Wrap 包装错误
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return err
}

// WrapMsg 包装错误并添加消息
func WrapMsg(err error, msg string, kv ...any) error {
	if err == nil {
		return nil
	}
	if codeErr, ok := err.(CodeError); ok {
		return codeErr.WrapMsg(msg, kv...)
	}
	return &errorWrapper{error: err, s: toString(msg, kv)}
}

// Unwrap returns the underlying error from a wrapped error
func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		inner := unwrap.Unwrap()
		if inner == nil {
			return err
		}
		err = inner
	}
	return err
}

// ErrPanic creates a CodeError from a panic value
func ErrPanic(r any) error {
	if r == nil {
		return nil
	}
	return &codeError{
		code:   ServerInternalError,
		msg:    "panic error",
		detail: fmt.Sprint(r),
	}
}

// errorWrapper wraps an error with additional message
type errorWrapper struct {
	error
	s string
}

func (e *errorWrapper) Error() string {
	return fmt.Sprintf("%s %s", e.error, e.s)
}

func (e *errorWrapper) Unwrap() error {
	return e.error
}

func (e *errorWrapper) Is(err error) bool {
	if err == nil {
		return false
	}
	var t *errorWrapper
	return errors.As(err, &t) && e.s == t.s
}

// toString converts message and key-value pairs to string
func toString(s string, kv []any) string {
	if len(kv) == 0 {
		return s
	}
	var buf bytes.Buffer
	buf.WriteString(s)
	for i := 0; i < len(kv); i += 2 {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		key := fmt.Sprintf("%v", kv[i])
		buf.WriteString(key)
		buf.WriteString("=")
		if i+1 < len(kv) {
			buf.WriteString(fmt.Sprintf("%v", kv[i+1]))
		} else {
			buf.WriteString("MISSING")
		}
	}
	return buf.String()
}

// Predefined errors (all error codes in one place, matching official OpenIM)
var (
	ErrSecretNotChanged = NewCodeError(SecretNotChangedError, "secret not changed, please change secret in config/share.yml for security reasons")

	ErrDatabase         = NewCodeError(DatabaseError, "DatabaseError")
	ErrNetwork          = NewCodeError(NetworkError, "NetworkError")
	ErrCallback         = NewCodeError(CallbackError, "CallbackError")
	ErrCallbackContinue = NewCodeError(CallbackError, "ErrCallbackContinue")

	ErrInternalServer = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrArgs           = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission   = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrDuplicateKey   = NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrRecordNotFound = NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrUserIDNotFound  = NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrGroupIDNotFound = NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrGroupIDExisted  = NewCodeError(GroupIDExisted, "GroupIDExisted")

	ErrNotInGroupYet       = NewCodeError(NotInGroupYetError, "NotInGroupYetError")
	ErrDismissedAlready    = NewCodeError(DismissedAlreadyError, "DismissedAlreadyError")
	ErrRegisteredAlready   = NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")
	ErrGroupTypeNotSupport = NewCodeError(GroupTypeNotSupport, "")
	ErrGroupRequestHandled = NewCodeError(GroupRequestHandled, "GroupRequestHandled")

	ErrData             = NewCodeError(DataError, "DataError")
	ErrTokenExpired     = NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid     = NewCodeError(TokenInvalidError, "TokenInvalidError")
	ErrTokenMalformed   = NewCodeError(TokenMalformedError, "TokenMalformedError")
	ErrTokenNotValidYet = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError")
	ErrTokenUnknown     = NewCodeError(TokenUnknownError, "TokenUnknownError")
	ErrTokenKicked      = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = NewCodeError(TokenNotExistError, "TokenNotExistError")

	ErrMessageHasReadDisable = NewCodeError(MessageHasReadDisable, "MessageHasReadDisable")

	ErrCanNotAddYourself   = NewCodeError(CanNotAddYourselfError, "CanNotAddYourselfError")
	ErrBlockedByPeer       = NewCodeError(BlockedByPeer, "BlockedByPeer")
	ErrNotPeersFriend      = NewCodeError(NotPeersFriend, "NotPeersFriend")
	ErrRelationshipAlready = NewCodeError(RelationshipAlreadyError, "RelationshipAlreadyError")

	ErrMutedInGroup     = NewCodeError(MutedInGroup, "MutedInGroup")
	ErrMutedGroup       = NewCodeError(MutedGroup, "MutedGroup")
	ErrMsgAlreadyRevoke = NewCodeError(MsgAlreadyRevoke, "MsgAlreadyRevoke")

	ErrConnOverMaxNumLimit = NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")

	ErrConnArgsErr          = NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
	ErrPushMsgErr           = NewCodeError(PushMsgErr, "push msg err")
	ErrIOSBackgroundPushErr = NewCodeError(IOSBackgroundPushErr, "ios background push err")

	ErrFileUploadedExpired = NewCodeError(FileUploadedExpiredError, "FileUploadedExpiredError")
)
