/// All error codes for fdim (frontend SDK + backend server).
/// Must stay in sync with backend: fdim/pkg/errs/errs.go
class SDKErrorCode {
  // ============================================================
  // Client SDK Error Codes (10000+, generated locally by client)
  // ============================================================

  /// Network Request Error
  static const int networkRequestError = 10000;

  /// Network Waiting Timeout Error
  static const int networkWaitTimeoutError = 10001;

  /// Parameter Error
  static const int parameterError = 10002;

  /// Context Timeout Error, usually when the user has already logged out
  static const int contextTimeoutError = 10003;

  /// Resources not loaded completely, usually uninitialized or login hasn't completed
  static const int resourceNotLoaded = 10004;

  /// Unknown Error, check the error message for details
  static const int unknownError = 10005;

  /// SDK Internal Error, check the error message for details
  static const int sdkInternalError = 10006;

  /// This user has set not to be added
  static const int refuseToAddFriends = 10013;

  /// User does not exist or is not registered
  static const int userNotExistOrNotRegistered = 10100;

  /// User has already logged out
  static const int userHasLoggedOut = 10101;

  /// User is attempting to log in again, check the login status to avoid duplicate logins
  static const int repeatLogin = 10102;

  /// The file to upload does not exist
  static const int uploadFileNotExist = 10200;

  /// Message decompression failed
  static const int messageDecompressionFailed = 10201;

  /// Message decoding failed
  static const int messageDecodingFailed = 10202;

  /// Unsupported long connection binary protocol
  static const int unsupportedLongConnection = 10203;

  /// Message sent multiple times
  static const int messageRepeated = 10204;

  /// Message content type not supported
  static const int messageContentTypeNotSupported = 10205;

  /// Unsupported session operation
  static const int unsupportedSessionOperation = 10301;

  /// Group ID does not exist (client-side check)
  static const int groupIDNotExist = 10400;

  /// Group type is incorrect (client-side check)
  static const int wrongGroupType = 10401;

  // ============================================================
  // Server Error Codes (returned by backend in errCode field)
  // Matches: fdim/pkg/errs/errs.go
  // ============================================================

  // --- General ---
  /// Server Internal Error
  static const int serverInternalError = 500;

  /// Input parameter error
  static const int argsError = 1001;

  /// Insufficient Permissions
  static const int noPermissionError = 1002;

  /// Duplicate Database Primary Key
  static const int duplicateKeyError = 1003;

  /// Record does not exist
  static const int recordNotFoundError = 1004;

  /// Secret not changed
  static const int secretNotChangedError = 1050;

  // --- Account ---
  /// UserID does not exist or is not registered
  static const int userIDNotFoundError = 1101;

  /// User is already registered
  static const int registeredAlreadyError = 1102;

  // --- Group ---
  /// GroupID does not exist
  static const int groupIDNotFoundError = 1201;

  /// GroupID already exists
  static const int groupIDExisted = 1202;

  /// Not in the group yet
  static const int notInGroupYetError = 1203;

  /// Group has already been dismissed
  static const int dismissedAlreadyError = 1204;

  /// Group type not supported
  static const int groupTypeNotSupport = 1205;

  /// Group application has already been processed
  static const int groupRequestHandled = 1206;

  // --- Relationship ---
  /// Cannot add yourself as a friend
  static const int canNotAddYourselfError = 1301;

  /// Blocked by the peer
  static const int blockedByPeer = 1302;

  /// Not the peer's friend
  static const int notPeersFriend = 1303;

  /// Already in a friend relationship
  static const int relationshipAlreadyError = 1304;

  // --- Message ---
  /// Message read function is turned off
  static const int messageHasReadDisable = 1401;

  /// Member muted in the group
  static const int mutedInGroup = 1402;

  /// Group is muted
  static const int mutedGroup = 1403;

  /// Message already revoked
  static const int msgAlreadyRevoke = 1404;

  /// License expired
  static const int licenseExpired = 1405;

  // --- Token ---
  /// Token has expired
  static const int tokenExpiredError = 1501;

  /// Invalid token
  static const int tokenInvalidError = 1502;

  /// Token format error (malformed)
  static const int tokenMalformedError = 1503;

  /// Token has not yet taken effect
  static const int tokenNotValidYetError = 1504;

  /// Unknown token error
  static const int tokenUnknownError = 1505;

  /// Token kicked (logged in from another device)
  static const int tokenKickedError = 1506;

  /// Token does not exist
  static const int tokenNotExistError = 1507;

  // --- Long Connection Gateway ---
  /// Number of connections exceeds gateway's maximum limit
  static const int connOverMaxNumLimit = 1601;

  /// Handshake parameter error
  static const int connArgsErr = 1602;

  /// Push message error
  static const int pushMsgErr = 1603;

  /// iOS background push error
  static const int iOSBackgroundPushErr = 1604;

  // --- S3 ---
  /// File upload expired
  static const int fileUploadedExpiredError = 1701;

  // --- Callback ---
  /// Callback error
  static const int callbackError = 80000;

  // --- Infrastructure ---
  /// Database error (redis/mysql, etc.)
  static const int databaseError = 90002;

  /// Network error (server-side)
  static const int serverNetworkError = 90004;

  /// Data error
  static const int dataError = 90007;

  // ============================================================
  // Chat/Account Server Error Codes (returned by chat-api)
  // Note: 10001-10014 range overlaps with SDK client codes above,
  // but these come from server HTTP responses, not local SDK.
  // Backend: fdim/pkg/errs/errs.go (FormattingError..InvitationError)
  // ============================================================

  /// Error in formatting
  static const int formattingError = 10001;

  /// User has already registered (chat layer)
  static const int chatHasRegistered = 10002;

  /// User is not registered (chat layer)
  static const int chatNotRegistered = 10003;

  /// Password error
  static const int passwordErr = 10004;

  /// Error in getting IM token
  static const int getIMTokenErr = 10005;

  /// Repeat sending verification code
  static const int repeatSendCode = 10006;

  /// Error in sending code via email
  static const int mailSendCodeErr = 10007;

  /// Error in sending code via SMS
  static const int smsSendCodeErr = 10008;

  /// Verification code is invalid or expired
  static const int codeInvalidOrExpired = 10009;

  /// Registration failed
  static const int registerFailed = 10010;

  /// Resetting password failed
  static const int resetPasswordFailed = 10011;

  /// Registration limit exceeded
  static const int registerLimit = 10012;

  /// Login limit exceeded
  static const int loginLimit = 10013;

  /// Invitation error
  static const int invitationError = 10014;

  // ============================================================
  // Backward compatibility aliases (old field names)
  // ============================================================
  static const int serverParameterError = argsError;
  static const int insufficientPermissions = noPermissionError;
  static const int duplicateDatabasePrimaryKey = duplicateKeyError;
  static const int databaseRecordNotFound = recordNotFoundError;
  static const int userIDNotExist = userIDNotFoundError;
  static const int userAlreadyRegistered = registeredAlreadyError;
  static const int groupNotExis = groupIDNotFoundError;
  static const int groupAlreadyExists = groupIDExisted;
  static const int userIsNotInGroup = notInGroupYetError;
  static const int groupDisbanded = dismissedAlreadyError;
  static const int groupApplicationHasBeenProcessed = groupRequestHandled;
  static const int notAddMyselfAsAFriend = canNotAddYourselfError;
  static const int hasBeenBlocked = blockedByPeer;
  static const int notFriend = notPeersFriend;
  static const int alreadyAFriendRelationship = relationshipAlreadyError;
  static const int messageReadFunctionIsTurnedOff = messageHasReadDisable;
  static const int youHaveBeenBanned = mutedInGroup;
  static const int groupHasBeenBanned = mutedGroup;
  static const int messageHasBeenRetracted = msgAlreadyRevoke;
  static const int tokenHasExpired = tokenExpiredError;
  static const int tokenInvalid = tokenInvalidError;
  static const int tokenFormatError = tokenMalformedError;
  static const int tokenHasNotYetTakenEffect = tokenNotValidYetError;
  static const int unknownTokenError = tokenUnknownError;
  static const int thekickedOutTokenIsInvalid = tokenKickedError;
  static const int tokenNotExist = tokenNotExistError;
  static const int connectionsExceedsMaximumLimit = connOverMaxNumLimit;
  static const int handshakeParameterError = connArgsErr;
  static const int fileUploadExpired = fileUploadedExpiredError;
}
