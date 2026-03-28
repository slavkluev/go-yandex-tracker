# Filters

The Filters service manages saved issue filters (queries).

## Usage

### List Filters

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

filters, _, err := client.Filters.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, f := range filters {
    fmt.Println(*f.Name)
}
```

### Create a Filter

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

filter, _, err := client.Filters.Create(context.Background(), &tracker.FilterRequest{
    Name:   tracker.Ptr("My Filter"),
    Filter: map[string]any{"queue": "QUEUE", "assignee": "me"},
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*filter.Name)
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all saved filters |
| `Get` | Get a filter by ID |
| `Create` | Create a saved filter |
| `Update` | Update a saved filter |
| `Delete` | Delete a saved filter |

## See Also

- [ExampleFiltersService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-FiltersService.List)
- [ExampleFiltersService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-FiltersService.Create)
- [Error Handling](errors.md)
