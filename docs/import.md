# Import

The Import service imports issues, comments, links, and files from external systems into Yandex Tracker.

## Usage

### Import an Issue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issue, _, err := client.Import.ImportIssue(context.Background(), &tracker.ImportIssueRequest{
    Queue:     tracker.Ptr("QUEUE"),
    Summary:   tracker.Ptr("Imported issue"),
    CreatedAt: tracker.Ptr("2024-01-15T10:00:00.000+0000"),
    CreatedBy: tracker.Ptr("admin"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*issue.Key)
```

### Import a Comment

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

comment, _, err := client.Import.ImportComment(context.Background(), "QUEUE-1", &tracker.ImportCommentRequest{
    Text:      tracker.Ptr("Imported comment"),
    CreatedAt: tracker.Ptr("2024-01-15T11:00:00.000+0000"),
    CreatedBy: tracker.Ptr("admin"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*comment.Text)
```

## Methods

| Method | Description |
|--------|-------------|
| `ImportIssue` | Import an issue with preserved timestamps |
| `ImportComment` | Import a comment to an issue |
| `ImportLink` | Import a link between issues |
| `ImportFile` | Import a file attachment to an issue |
| `ImportCommentFile` | Import a file attachment to a comment |

## See Also

- [ExampleImportService_ImportIssue](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ImportService.ImportIssue)
- [Error Handling](errors.md)
