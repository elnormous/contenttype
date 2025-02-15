package contenttype

import "errors"

var (
	// ErrInvalidMediaType is returned when the media type in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidMediaType = errors.New("invalid media type")
	// ErrInvalidMediaRange is returned when the range of media types in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidMediaRange = errors.New("invalid media range")
	// ErrInvalidParameter is returned when the media type parameter in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidParameter = errors.New("invalid parameter")
	// ErrInvalidExtensionParameter is returned when the media type extension parameter in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidExtensionParameter = errors.New("invalid extension parameter")
	// ErrNoAcceptableTypeFound is returned when Accept header contains only media types that are not in the acceptable media type list.
	ErrNoAcceptableTypeFound = errors.New("no acceptable type found")
	// ErrNoAvailableTypeGiven is returned when the acceptable media type list is empty.
	ErrNoAvailableTypeGiven = errors.New("no available type given")
	// ErrInvalidWeight is returned when the media type weight in Accept header is syntactically invalid.
	ErrInvalidWeight = errors.New("invalid weight")
	// ErrInvalidLanguage is returned when the language is syntactically invalid.
	ErrInvalidLanguage = errors.New("invalid language")
)
