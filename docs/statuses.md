# Statuses

The Statuses service lists available workflow statuses in the Yandex Tracker organization.

## Usage

### List Statuses

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

statuses, _, err := client.Statuses.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, s := range statuses {
    fmt.Println(*s.Key, *s.Name)
}
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all workflow statuses |

## See Also

- [ExampleStatusesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-StatusesService.List)
- [Error Handling](errors.md)
