# Issue Types

The Issue Types service lists available issue types in the Yandex Tracker organization.

## Usage

### List Issue Types

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

for t, err := range client.IssueTypes.ListIter(context.Background(), nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*t.Key, *t.Name)
}
```

`List` returns one page: `&tracker.ListOptions{Page: 2}` chooses which and `PerPage` its size, while `nil` gets the first page at the server's default size.

## Methods

| Method | Description |
|--------|-------------|
| `List` | List one page of issue types |
| `ListIter` | Iterate over the issue types of every page |

## See Also

- [ExampleIssueTypesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssueTypesService.List)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
