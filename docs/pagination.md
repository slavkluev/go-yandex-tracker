# Pagination

A list method returns one page. Issue search, the queue, user, status, priority, resolution and issue type lists, issue comments and the issue changelog also have an iterator method, named after the method with an `Iter` suffix, that walks every page and yields every item. The entity search, comment and event lists (`Entities.Search`, `Entities.ListComments`, `Entities.GetEvents`) have none. Issue search also supports scroll-based pagination for large result sets. Pagination metadata is returned in the `Response` struct.

## Iterators

An iterator takes the same arguments as the method it pages and returns an `iter.Seq2[T, error]`:

```go
req := &tracker.IssueSearchRequest{
    Filter: map[string]any{"queue": "QUEUE"},
}
opts := &tracker.IssueSearchOptions{
    ListOptions: tracker.ListOptions{PerPage: 100},
}

for issue, err := range client.Issues.SearchIter(ctx, req, opts) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*issue.Key)
}
```

| Iterator | Next page | Last page |
|----------|-----------|-----------|
| `Issues.SearchIter`, `Queues.ListIter`, `Users.ListIter`, `Statuses.ListIter`, `Priorities.ListIter`, `Resolutions.ListIter`, `IssueTypes.ListIter` | `page` + 1 | page `X-Total-Pages`, or the first empty page when that header is missing |
| `Issues.GetChangelogIter`, `Issues.ListCommentsIter` | `id` set to the ID of the previous page's last item | the first empty page |

- The first request is the one the method sends for the same arguments, so `Page` or `ID` in the options sets where the iteration starts, and `PerPage` sets the page size.
- The iterator copies the options when it is called and never changes them.
- Tracker sends `rel="next"` from the last full page of a cursor list too, so a cursor iteration ends one request after its last item, on the empty page.
- An error is yielded once, with a nil item, and ends the iteration. A cursor page whose last item has no ID or repeats the cursor yields an error instead of requesting the same page again.
- Breaking out of the loop sends no further request.

## Page-Based Pagination

Embed `ListOptions` in the request options struct to choose one page and its size:

```go
opts := &tracker.IssueSearchOptions{
    ListOptions: tracker.ListOptions{Page: 2, PerPage: 50},
}

issues, resp, err := client.Issues.Search(ctx, req, opts)
if err != nil {
    log.Fatal(err)
}

for _, issue := range issues {
    fmt.Println(*issue.Key)
}
```

The `Response` includes pagination headers:

```go
fmt.Println("Total results:", resp.TotalCount)
fmt.Println("Total pages:", resp.TotalPages)
```

## Cursor-Based Pagination

Issue comments and the issue changelog page by item ID: set `ID` in `CommentListOptions` or `ChangelogOptions` to the ID of the last item you have, and the next page starts after it.

```go
comments, _, err := client.Issues.ListComments(ctx, "QUEUE-1", &tracker.CommentListOptions{
    ID:      lastID,
    PerPage: 50,
})
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

- [ExampleIssuesService_SearchIter](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.SearchIter)
- [ExampleIssuesService_Search_pagination](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search-pagination)
- [ExampleIssuesService_Search](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search)
- [Authentication](auth.md)
- [Error Handling](errors.md)
