# Issue Types

The Issue Types service lists available issue types in the Yandex Tracker organization.

## Usage

### List Issue Types

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issueTypes, _, err := client.IssueTypes.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, t := range issueTypes {
    fmt.Println(*t.Key, *t.Name)
}
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all issue types |

## See Also

- [ExampleIssueTypesService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssueTypesService.List)
- [Error Handling](errors.md)
