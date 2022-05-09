package exception

type ValidationError struct {
	Message string
}

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
func (validationError ValidationError) Error() string {

	return validationError.Message
}
