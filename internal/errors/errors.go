// Package errors provides a custom error type that can be used to create
// named errors
package errors

// Error provides a custom type for creating named errors in various packages
type Error string

// Error implements the Error interface
func (e Error) Error() string { return string(e) }
