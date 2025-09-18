package quality

var (
	ErrTopicExists      = NewError("topic already exists")
	ErrTopicNotFound    = NewError("topic not found")
	ErrSubscriberExists = NewError("subscriber already exists for this topic")

	ErrInvalidPayload = NewError("")
	ErrPayloadBinding = NewError("")
	ErrInvalidAction = NewError("")

	// testing & quality
	ErrNotImplemented   = NewError("not implemented")

	// http 401
	ErrUnauthorized     = NewError("not authorized")
	// http 403
	ErrForbidden        = NewError("forbidden")
	// http 404
	ErrNotFound         = NewError("not found")
	// http 500
	ErrInternal         = NewError("internal error")

	// generic execution error
	ErrExecution        = NewError("execution error")

	// filesystem error
	ErrFsIO        = NewError("filesystem error")

	// password validation error
	ErrPasswordNotMatch = NewError("password does not match")
	// generic validation error
	ErrValidation       = NewError("validation error")

	// data integrity errors
	ErrEmailExists      = NewError("email exists")
	ErrDuplicateEntry   = NewError("unique key")
	ErrRecordNotFound   = NewError("record not found")

	// DB errors
	ErrDsnNotSet        = NewError("dsn not set")
	ErrSQLConnection    = NewError("sql conn err")
	ErrSQLQuery         = NewError("sql query err")
)
