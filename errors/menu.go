package errors

var (
	MenuRecordNotFound          = New("menu record not found")
	MenuAlreadyExists           = New("menu already exists")
	MenuInvalidParent           = New("menu invalid parent")
	MenuNotAllowDeleteWithChild = New("contains children, cannot be deleted")
)
