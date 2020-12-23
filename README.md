# Content-Type support library for Go

## Usage

```go
import (
	"log"
	"github.com/elnormous/contenttype"
)

func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
    mediaType, parameters, err := contenttype.GetMediaType(request)
    if err != nil {
        // handle the error
    }
    log.Println("Media type:", result, "parameters:", parameters)
}
```