# Priorities

The Priorities service lists available issue priorities in the Yandex Tracker organization.

## Usage

### List Priorities

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

priorities, _, err := client.Priorities.List(context.Background(), nil)
if err != nil {
    log.Fatal(err)
}

for _, p := range priorities {
    fmt.Println(*p.Key, *p.Name)
}
```

Note: `List` accepts an optional `*PriorityListOptions` parameter (pass `nil` for defaults).

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all issue priorities |

## See Also

- [ExamplePrioritiesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-PrioritiesService.List)
- [Error Handling](errors.md)
