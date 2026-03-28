# Error Handling

API errors are returned as `*tracker.ErrorResponse`, which implements the `error` interface. The library provides helper functions for common HTTP status codes and supports Go's `errors.As` for detailed inspection.

## Helper Functions

Check common error types without manual status code inspection:

```go
issue, _, err := client.Issues.Get(ctx, "QUEUE-999", nil)
if tracker.IsNotFound(err) {
    log.Println("issue not found")
    return
}
if tracker.IsForbidden(err) {
    log.Println("access denied")
    return
}
if tracker.IsRateLimited(err) {
    log.Println("rate limited, retry later")
    return
}
if err != nil {
    log.Fatal(err)
}
```

## Inspecting ErrorResponse

Use `errors.As` to access the full error details, including status code, error messages, and field-level errors:

```go
issue, _, err := client.Issues.Get(ctx, "QUEUE-999", nil)
if err != nil {
    var errResp *tracker.ErrorResponse
    if errors.As(err, &errResp) {
        fmt.Println("Status:", errResp.Response.StatusCode)
        fmt.Println("Messages:", errResp.ErrorMessages)
        fmt.Println("Errors:", errResp.Errors)
    }
    log.Fatal(err)
}
```

## See Also

- [ExampleIssuesService_Get](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Get)
- [Authentication](auth.md)
- [Pagination](pagination.md)
