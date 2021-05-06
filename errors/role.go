package errors

var (
	RoleRecordNotFound         = New("role record not found")
	RoleIsDisable              = New("role is disabled")
	RoleAlreadyExists          = New("role already exists")
	RoleNotAllowDeleteWithUser = New("used by users, cannot be deleted")
)
