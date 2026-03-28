# External Links

The External Links service manages links between Yandex Tracker issues and external applications.

## Usage

### List External Links

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

links, _, err := client.ExternalLinks.ListLinks(context.Background(), "QUEUE-1")
if err != nil {
    log.Fatal(err)
}

for _, link := range links {
    fmt.Println(*link.ID)
}
```

### Create an External Link

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

link, _, err := client.ExternalLinks.CreateLink(context.Background(), "QUEUE-1", &tracker.ExternalLinkCreateRequest{
    Relationship: tracker.Ptr("relates"),
    Key:          tracker.Ptr("EXT-123"),
}, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*link.ID)
```

## Methods

| Method | Description |
|--------|-------------|
| `ListApplications` | List registered external applications |
| `ListLinks` | List external links on an issue |
| `CreateLink` | Create an external link on an issue |
| `DeleteLink` | Delete an external link |

## See Also

- [ExampleExternalLinksService_ListLinks](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ExternalLinksService.ListLinks)
- [Error Handling](errors.md)
