package internal

var (
	NotFoundErr   = NotFoundError{BaseError{Message: "not found"}}
	InvalidBinErr = InvalidBinError{BaseError{Message: "bid bin is empty or not exists"}}
)

type (
	BaseError struct {
		Message string
	}

	NotFoundError struct {
		BaseError
	}

	InvalidBinError struct {
		BaseError
	}
)

func (e BaseError) Error() string {
	return e.Message
}
