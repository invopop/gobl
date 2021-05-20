package dsig

// Error defines the standard error messages supported by this
// JWS library.
type Error string

// Standard error messages
const (
	ErrKeyPublic    Error = "cannot sign with public key"
	ErrKeyInvalid   Error = "key is not valid"
	ErrKeyMismatch  Error = "key mismatch"
	ErrVerifyFailed Error = "verification failed"
)

// Error provides the standard error response text.
func (e Error) Error() string {
	return string(e)
}
