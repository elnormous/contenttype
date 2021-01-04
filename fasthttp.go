package contenttype

import (
	"github.com/valyala/fasthttp"
)

// Gets the content of Content-Type header, parses it, and returns the parsed MediaType
// If the request does not contain the Content-Type header, an empty MediaType is returned
func GetMediaTypeFastHTTP(request *fasthttp.Request) (MediaType, error) {
	// RFC 7231, 3.1.1.5. Content-Type
	contentTypeHeaders := string(request.Header.Peek("Content-Type"))
	if len(contentTypeHeaders) == 0 {
		return MediaType{}, nil
	}

	s := contentTypeHeaders
	mediaType := MediaType{}
	var consumed bool
	mediaType.Type, mediaType.Subtype, s, consumed = consumeType(s)
	if !consumed {
		return MediaType{}, ErrInvalidMediaType
	}

	mediaType.Parameters = make(Parameters)

	for len(s) > 0 && s[0] == ';' {
		s = s[1:] // skip the semicolon

		key, value, remaining, consumed := consumeParameter(s)
		if !consumed {
			return MediaType{}, ErrInvalidParameter
		}

		s = remaining

		mediaType.Parameters[key] = value
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return MediaType{}, ErrInvalidMediaType
	}

	return mediaType, nil
}

// Choses a media type from available media types according to the Accept
// Returns the most suitable media type or an error if no type can be selected
func GetAcceptableMediaTypeFastHTTP(request *fasthttp.Request, availableMediaTypes []MediaType) (MediaType, Parameters, error) {
	// RFC 7231, 5.3.2. Accept
	if len(availableMediaTypes) == 0 {
		return MediaType{}, Parameters{}, ErrNoAvailableTypeGiven
	}

	acceptHeaders := string(request.Header.Peek("Accept"))
	if len(acceptHeaders) == 0 {
		return availableMediaTypes[0], Parameters{}, nil
	}

	s := acceptHeaders

	weights := make([]struct {
		mediaType           MediaType
		extensionParameters Parameters
		weight              int
		order               int
	}, len(availableMediaTypes))

	for mediaTypeCount := 0; len(s) > 0; mediaTypeCount++ {
		if mediaTypeCount > 0 {
			// every media type after the first one must start with a comma
			if s[0] != ',' {
				break
			}
			s = s[1:] // skip the comma
		}

		acceptableMediaType := MediaType{}
		var consumed bool
		acceptableMediaType.Type, acceptableMediaType.Subtype, s, consumed = consumeType(s)
		if !consumed {
			return MediaType{}, Parameters{}, ErrInvalidMediaType
		}

		acceptableMediaType.Parameters = make(Parameters)
		weight := 1000 // 1.000

		// media type parameters
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			var key, value string
			key, value, s, consumed = consumeParameter(s)
			if !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			if key == "q" {
				weight, consumed = getWeight(value)
				if !consumed {
					return MediaType{}, Parameters{}, ErrInvalidWeight
				}
				break // "q" parameter separates media type parameters from Accept extension parameters
			}

			acceptableMediaType.Parameters[key] = value
		}

		extensionParameters := make(Parameters)
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			key, value, remaining, consumed := consumeParameter(s)
			if !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			s = remaining

			extensionParameters[key] = value
		}

		for i := 0; i < len(availableMediaTypes); i++ {
			if compareMediaTypes(acceptableMediaType, availableMediaTypes[i]) &&
				getPrecedence(acceptableMediaType, weights[i].mediaType) {
				weights[i].mediaType = acceptableMediaType
				weights[i].extensionParameters = extensionParameters
				weights[i].weight = weight
				weights[i].order = mediaTypeCount
			}
		}

		s = skipWhiteSpaces(s)
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return MediaType{}, Parameters{}, ErrInvalidMediaRange
	}

	resultIndex := -1
	for i := 0; i < len(availableMediaTypes); i++ {
		if resultIndex != -1 {
			if weights[i].weight > weights[resultIndex].weight ||
				(weights[i].weight == weights[resultIndex].weight && weights[i].order < weights[resultIndex].order) {
				resultIndex = i
			}
		} else if weights[i].weight > 0 {
			resultIndex = i
		}
	}

	if resultIndex == -1 {
		return MediaType{}, Parameters{}, ErrNoAcceptableTypeFound
	}

	return availableMediaTypes[resultIndex], weights[resultIndex].extensionParameters, nil
}
