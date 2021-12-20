package provider

// Error provides a custom type for creating named errors in various packages
type Error string

// Error implements the Error interface
func (e Error) Error() string { return string(e) }
