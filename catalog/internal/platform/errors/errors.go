package errors

// Error is a handy type alias that lets us create constant error messages
type Error string

// Error is the "stringer" function, which completes our error interface implementation.
func (e Error) Error() string {
	return string(e)
}

const (
	// BadSearchTerm indicates a search term validation failure
	BadSearchTerm = Error("Bad search term")

	// NoSuchCategory indicates a request for a category that doesn't exist.
	NoSuchCategory = Error("No such category")

	// NoSuchProduct indicates a request for a non-existent product
	NoSuchProduct = Error("No such product")
)
