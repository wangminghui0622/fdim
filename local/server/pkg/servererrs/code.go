// Deprecated: All error codes are now in pkg/errs/errs.go.
// This file re-exports them for backward compatibility.
package servererrs

import "fdim/pkg/errs"

const (
	UnknownCode          = errs.UnknownCode
	FormattingError      = errs.FormattingError
	HasRegistered        = errs.HasRegistered
	NotRegistered        = errs.NotRegistered
	PasswordErr          = errs.PasswordErr
	GetIMTokenErr        = errs.GetIMTokenErr
	RepeatSendCode       = errs.RepeatSendCode
	MailSendCodeErr      = errs.MailSendCodeErr
	SmsSendCodeErr       = errs.SmsSendCodeErr
	CodeInvalidOrExpired = errs.CodeInvalidOrExpired
	RegisterFailed       = errs.RegisterFailed
	ResetPasswordFailed  = errs.ResetPasswordFailed
	RegisterLimit        = errs.RegisterLimit
	LoginLimit           = errs.LoginLimit
	InvitationError      = errs.InvitationError
)

const (
	NoError               = errs.NoError
	DatabaseError         = errs.DatabaseError
	NetworkError          = errs.NetworkError
	DataError             = errs.DataError
	CallbackError         = errs.CallbackError
	ServerInternalError   = errs.ServerInternalError
	ArgsError             = errs.ArgsError
	NoPermissionError     = errs.NoPermissionError
	DuplicateKeyError     = errs.DuplicateKeyError
	RecordNotFoundError   = errs.RecordNotFoundError
	SecretNotChangedError = errs.SecretNotChangedError

	UserIDNotFoundError    = errs.UserIDNotFoundError
	RegisteredAlreadyError = errs.RegisteredAlreadyError

	GroupIDNotFoundError  = errs.GroupIDNotFoundError
	GroupIDExisted        = errs.GroupIDExisted
	NotInGroupYetError    = errs.NotInGroupYetError
	DismissedAlreadyError = errs.DismissedAlreadyError
	GroupTypeNotSupport   = errs.GroupTypeNotSupport
	GroupRequestHandled   = errs.GroupRequestHandled

	CanNotAddYourselfError   = errs.CanNotAddYourselfError
	BlockedByPeer            = errs.BlockedByPeer
	NotPeersFriend           = errs.NotPeersFriend
	RelationshipAlreadyError = errs.RelationshipAlreadyError

	MessageHasReadDisable = errs.MessageHasReadDisable
	MutedInGroup          = errs.MutedInGroup
	MutedGroup            = errs.MutedGroup
	MsgAlreadyRevoke      = errs.MsgAlreadyRevoke

	TokenExpiredError     = errs.TokenExpiredError
	TokenInvalidError     = errs.TokenInvalidError
	TokenMalformedError   = errs.TokenMalformedError
	TokenNotValidYetError = errs.TokenNotValidYetError
	TokenUnknownError     = errs.TokenUnknownError
	TokenKickedError      = errs.TokenKickedError
	TokenNotExistError    = errs.TokenNotExistError

	ConnOverMaxNumLimit  = errs.ConnOverMaxNumLimit
	ConnArgsErr          = errs.ConnArgsErr
	PushMsgErr           = errs.PushMsgErr
	IOSBackgroundPushErr = errs.IOSBackgroundPushErr

	FileUploadedExpiredError = errs.FileUploadedExpiredError
)
