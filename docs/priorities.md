# Priorities

The Priorities service lists available issue priorities in the Yandex Tracker organization.

## Usage

### List Priorities

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

for p, err := range client.Priorities.ListIter(context.Background(), nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*p.Key, *p.Name)
}
```

`List` returns one page. `*PriorityListOptions` chooses it through the embedded `ListOptions` and sets `Localized`; `nil` gets the first page at the server's default size.

## Methods

| Method | Description |
|--------|-------------|
| `List` | List one page of issue priorities |
| `ListIter` | Iterate over the issue priorities of every page |

## See Also

- [ExamplePrioritiesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-PrioritiesService.List)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
