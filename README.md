# Request ID middleware

## Usage

Just wrap handler with provided middleware:
```go                                                                      
handler = requestid.Middleware(handler, requestid.DefaultRequestIdProvider)
```
