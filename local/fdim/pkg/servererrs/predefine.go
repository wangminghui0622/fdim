// Deprecated: All predefined errors are now in pkg/errs/errs.go.
// This file re-exports them for backward compatibility.
package servererrs

import "fdim/pkg/errs"

var (
	ErrSecretNotChanged = errs.ErrSecretNotChanged
	ErrDatabase         = errs.ErrDatabase
	ErrNetwork          = errs.ErrNetwork
	ErrCallback         = errs.ErrCallback
	ErrCallbackContinue = errs.ErrCallbackContinue

	ErrInternalServer = errs.ErrInternalServer
	ErrArgs           = errs.ErrArgs
	ErrNoPermission   = errs.ErrNoPermission
	ErrDuplicateKey   = errs.ErrDuplicateKey
	ErrRecordNotFound = errs.ErrRecordNotFound

	ErrUserIDNotFound  = errs.ErrUserIDNotFound
	ErrGroupIDNotFound = errs.ErrGroupIDNotFound
	ErrGroupIDExisted  = errs.ErrGroupIDExisted

	ErrNotInGroupYet       = errs.ErrNotInGroupYet
	ErrDismissedAlready    = errs.ErrDismissedAlready
	ErrRegisteredAlready   = errs.ErrRegisteredAlready
	ErrGroupTypeNotSupport = errs.ErrGroupTypeNotSupport
	ErrGroupRequestHandled = errs.ErrGroupRequestHandled

	ErrData             = errs.ErrData
	ErrTokenExpired     = errs.ErrTokenExpired
	ErrTokenInvalid     = errs.ErrTokenInvalid
	ErrTokenMalformed   = errs.ErrTokenMalformed
	ErrTokenNotValidYet = errs.ErrTokenNotValidYet
	ErrTokenUnknown     = errs.ErrTokenUnknown
	ErrTokenKicked      = errs.ErrTokenKicked
	ErrTokenNotExist    = errs.ErrTokenNotExist

	ErrMessageHasReadDisable = errs.ErrMessageHasReadDisable

	ErrCanNotAddYourself   = errs.ErrCanNotAddYourself
	ErrBlockedByPeer       = errs.ErrBlockedByPeer
	ErrNotPeersFriend      = errs.ErrNotPeersFriend
	ErrRelationshipAlready = errs.ErrRelationshipAlready

	ErrMutedInGroup     = errs.ErrMutedInGroup
	ErrMutedGroup       = errs.ErrMutedGroup
	ErrMsgAlreadyRevoke = errs.ErrMsgAlreadyRevoke

	ErrConnOverMaxNumLimit = errs.ErrConnOverMaxNumLimit

	ErrConnArgsErr          = errs.ErrConnArgsErr
	ErrPushMsgErr           = errs.ErrPushMsgErr
	ErrIOSBackgroundPushErr = errs.ErrIOSBackgroundPushErr

	ErrFileUploadedExpired = errs.ErrFileUploadedExpired
)
