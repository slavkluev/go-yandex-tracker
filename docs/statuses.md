# Statuses

The Statuses service lists available workflow statuses in the Yandex Tracker organization.

## Usage

### List Statuses

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

for s, err := range client.Statuses.ListIter(context.Background(), nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*s.Key, *s.Name)
}
```

`List` returns one page: `&tracker.ListOptions{Page: 2}` chooses which and `PerPage` its size, while `nil` gets the first page at the server's default size.

## Methods

| Method | Description |
|--------|-------------|
| `List` | List one page of workflow statuses |
| `ListIter` | Iterate over the workflow statuses of every page |

## See Also

- [ExampleStatusesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-StatusesService.List)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
