# Resolutions

The Resolutions service lists available issue resolutions in the Yandex Tracker organization.

## Usage

### List Resolutions

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

resolutions, _, err := client.Resolutions.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, r := range resolutions {
    fmt.Println(*r.Key, *r.Name)
}
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all issue resolutions |

## See Also

- [ExampleResolutionsService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ResolutionsService.List)
- [Error Handling](errors.md)
