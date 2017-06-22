package errors

// Error is a handy type alias that lets us create constant error messages
type Error string

// Error is the "stringer" function, which completes our error interface implementation.
func (e Error) Error() string {
	return string(e)
}

const (
	// NoSuchSKU indicates a request for a non-existent or invalid SKU
	NoSuchSKU = Error("No such SKU")

	// NoSuchOrder indicates a request for a non-existent or invalid order ID
	NoSuchOrder = Error("No such Order")
)
