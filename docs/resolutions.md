# Resolutions

The Resolutions service lists available issue resolutions in the Yandex Tracker organization.

## Usage

### List Resolutions

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

for r, err := range client.Resolutions.ListIter(context.Background(), nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*r.Key, *r.Name)
}
```

`List` returns one page: `&tracker.ListOptions{Page: 2}` chooses which and `PerPage` its size, while `nil` gets the first page at the server's default size.

## Methods

| Method | Description |
|--------|-------------|
| `List` | List one page of resolutions |
| `ListIter` | Iterate over the resolutions of every page |

## See Also

- [ExampleResolutionsService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ResolutionsService.List)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
