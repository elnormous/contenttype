# Content-Type support library for Go

This library can be used to parse the value Content-Type header (if one is present) and select an acceptable media type from the Accept header of HTTP request.

## Usage

Media types are stored in `MediaType` structure which has `Type` (e.g. `application`), Subtype (e.g. `json`) and Parameters (e.g. `charset: utf-8`) attributes. Media types are not stored in a string because media type parameters are part of the media type ([RFC 7231, 3.1.1.1. Media Type](https://tools.ietf.org/html/rfc7231#section-3.1.1.1)). To convert a string to `MediaType` use `NewMediaType`. To convert `MediaType` back to string use `String` function. If the `Content-Type` header is not present in the request, an empty `MediaType` is returned.

To get the `MediaType` of the incoming request call `GetMediaType` and pass the `http.Request` pointer to it. The function will return error if the `Content-Type` header is malformed according to [RFC 7231, 3.1.1.5. Content-Type](https://tools.ietf.org/html/rfc7231#section-3.1.1.5).

To get an acceptable media type from an `Accept`  header of the incoming request call `GetAcceptableMediaType` and pass the `http.Request` pointer to it and an array of all the acceptable media types. The function will return the best match following the negotiation rules written in [RFC 7231, 5.3.2. Accept](https://tools.ietf.org/html/rfc7231#section-5.3.2) or an error if the header is malformed or the content type in the `Accept` header is not supported. If the `Accept` header is not present in the request, the first media type from the acceptable type list is returned.

```go
import (
	"log"
	"github.com/elnormous/contenttype"
)

func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
    mediaType, mediaTypeError := contenttype.GetMediaType(request)
    if mediaTypeError != nil {
        // handle the error
    }
    log.Println("Media type:", mediaType.String())

    availableMediaTypes := []MediaType{
        contenttype.NewMediaType("application/json"),
        contenttype.NewMediaType("application/xml"),
    }

    accepted, extParameters, acceptError := contenttype.GetAcceptableMediaType(request, availableMediaTypes)
    if acceptError != nil {
        // handle the error
    }
    log.Println("Accepted media type:", accepted.String(), "extension parameters:", extParameters)
}
```