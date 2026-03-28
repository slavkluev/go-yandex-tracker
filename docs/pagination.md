# Pagination

List methods support page-based pagination through `ListOptions`. Issue search additionally supports scroll-based pagination for large result sets. Pagination metadata is returned in the `Response` struct.

## Page-Based Pagination

Embed `ListOptions` in the request options struct to control page number and page size:

```go
opts := &tracker.IssueSearchOptions{
    ListOptions: tracker.ListOptions{Page: 1, PerPage: 50},
}

for {
    issues, resp, err := client.Issues.Search(ctx, req, opts)
    if err != nil {
        log.Fatal(err)
    }

    for _, issue := range issues {
        fmt.Println(*issue.Key)
    }

    if len(issues) < opts.PerPage {
        break
    }
    opts.Page++
}
```

The `Response` includes pagination headers:

```go
fmt.Println("Total results:", resp.TotalCount)
fmt.Println("Total pages:", resp.TotalPages)
```

## Scroll-Based Pagination

For large result sets (10,000+ issues), use scroll-based pagination with `ScrollSearch` and `ScrollNext`:

```go
// Start scroll search
issues, resp, err := client.Issues.ScrollSearch(ctx, req, &tracker.ScrollSearchOptions{
    ScrollType: "sorted",
    PerScroll:  100,
})
if err != nil {
    log.Fatal(err)
}

// Process first batch
for _, issue := range issues {
    fmt.Println(*issue.Key)
}

// Continue with scroll token
for resp.ScrollToken != "" {
    issues, resp, err = client.Issues.ScrollNext(ctx, req, resp.ScrollID, resp.ScrollToken)
    if err != nil {
        log.Fatal(err)
    }

    for _, issue := range issues {
        fmt.Println(*issue.Key)
    }
}
```

## Response Metadata

Every API response includes a `*tracker.Response` with pagination and rate limiting fields:

| Field | Type | Source Header | Description |
|-------|------|---------------|-------------|
| `TotalCount` | `int` | `X-Total-Count` | Total number of results |
| `TotalPages` | `int` | `X-Total-Pages` | Total number of pages |
| `ScrollID` | `string` | `X-Scroll-Id` | Scroll cursor ID |
| `ScrollToken` | `string` | `X-Scroll-Token` | Scroll continuation token |
| `RetryAfter` | `int` | `Retry-After` | Seconds to wait on 429 |

## See Also

- [ExampleIssuesService_Search_pagination](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search-pagination)
- [ExampleIssuesService_Search](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search)
- [Authentication](auth.md)
- [Error Handling](errors.md)
