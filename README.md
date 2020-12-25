# Content-Type support library for Go

## Usage

```go
import (
	"log"
	"github.com/elnormous/contenttype"
)

func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
    mediaType, err := contenttype.GetMediaType(request)
    if err != nil {
        // handle the error
    }
    log.Println("Media type:", result.Type, result.Subtype, "parameters:", result.Parameters)
}
```