package errors

var (
	UserRecordNotFound   = New("user record not found")
	UserInvalidPassword  = New("invalid user password")
	UserIsDisable        = New("user is disabled")
	UserPasswordRequired = New("user password is required")
	UserInvalidUsername  = New("invalid username")
	UserAlreadyExists    = New("user already exists")
	UserNoPermission     = New("user no permission")
)
