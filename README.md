# Content-Type support library for Go

## Usage

```go
import (
	"log"
	"github.com/elnormous/contenttype"
)

func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
    mediaType, mediaTypeError := contenttype.GetMediaType(request)
    if contentTypeError != nil {
        // handle the error
    }
    log.Println("Media type:", result.Type, result.Subtype, "parameters:", result.Parameters)

    availableMediaTypes := []MediaType{
        contenttype.NewMediaType("application/json"),
        contenttype.NewMediaType("application/xml"),
    }

    accepted, extParameters, acceptError := contenttype.GetAcceptableMediaType(request, availableMediaTypes)
    if acceptError != nil {
        // handle the error
    }
    log.Println("Accepted media type:", accepted.Type, accepted.Subtype, "parameters:", accepted.Parameters)
    log.Println("Extension parameters:", extParameters)
}
```