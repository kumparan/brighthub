package brighthub

import "errors"

var (
	// ErrUnauthorized :nodoc:
	ErrUnauthorized = errors.New("unauthorized")
	// ErrIllegalField :nodoc:
	ErrIllegalField = errors.New("spelling error or other use of non-existent field")
	// ErrMethodNotAllowed :nodoc:
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrDuplicateReferenceID :nodoc:
	ErrDuplicateReferenceID = errors.New("duplicate reference id")
	// ErrTooManyRequest :nodoc:
	ErrTooManyRequest = errors.New("too many request")
	// ErrDynamicDeliveryNotAllowed :nodoc:
	ErrDynamicDeliveryNotAllowed = errors.New("this account is not enabled for Dynamic Delivery, but a Dynamic Delivery profile was specified")
	// ErrRateLimitExceeded :nodoc:
	ErrRateLimitExceeded = errors.New("dynamic ingest job not created. reduce the number of concurrent jobs for this account before trying again")
	// ErrInternalError :nodoc:
	ErrInternalError = errors.New("internal error, please try again later")
	// ErrBadRequest :nodoc:
	ErrBadRequest = errors.New("unable to parse request body")
	// ErrResourceNotFound :nodoc:
	ErrResourceNotFound = errors.New("the api could not find the resource you requested")
	// ErrNotAvailable :nodoc:
	ErrNotAvailable = errors.New("the resource you are requesting is temporarily unavailable")
	// ErrProfileError :nodoc:
	ErrProfileError = errors.New("profile rendition count exceeds configured rendition limit")
)
