# Issues

The Issues service manages Yandex Tracker issues: create, read, update, search, and manage sub-resources (comments, attachments, worklogs, links, transitions, checklists).

## Usage

### Get an Issue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issue, _, err := client.Issues.Get(context.Background(), "QUEUE-1", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*issue.Key, *issue.Summary)
```

### Create an Issue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issue, _, err := client.Issues.Create(context.Background(), &tracker.IssueRequest{
    Summary: tracker.Ptr("Bug: login fails"),
    Queue:   tracker.Ptr("QUEUE"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*issue.Key)
```

### Search Issues

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issues, _, err := client.Issues.Search(context.Background(), &tracker.IssueSearchRequest{
    Filter: map[string]any{
        "queue":    "QUEUE",
        "assignee": "me",
    },
}, nil)
if err != nil {
    log.Fatal(err)
}

for _, issue := range issues {
    fmt.Println(*issue.Key)
}
```

### Edit an Issue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

issue, _, err := client.Issues.Edit(context.Background(), "QUEUE-1", &tracker.IssueRequest{
    Summary: tracker.Ptr("Updated summary"),
}, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*issue.Key, *issue.Summary)
```

### List Comments

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

comments, _, err := client.Issues.ListComments(context.Background(), "QUEUE-1", nil)
if err != nil {
    log.Fatal(err)
}

for _, comment := range comments {
    fmt.Println(*comment.Text)
}
```

## Methods

| Method | Description |
|--------|-------------|
| **Core** | |
| `Create` | Create a new issue |
| `Get` | Get an issue by key |
| `Edit` | Update an issue |
| `Search` | Search issues with filters |
| `ScrollSearch` | Start scroll-based search for large result sets |
| `ScrollNext` | Continue scroll-based search with token |
| `Count` | Count issues matching a filter |
| `Move` | Move an issue to another queue |
| **Comments** | |
| `ListComments` | List comments on an issue |
| `CreateComment` | Add a comment to an issue |
| `EditComment` | Edit an issue comment |
| `DeleteComment` | Delete an issue comment |
| **Attachments** | |
| `ListAttachments` | List attachments on an issue |
| `UploadAttachment` | Upload a file attachment to an issue |
| `DeleteAttachment` | Delete an attachment |
| `UploadTempFile` | Upload a temporary file |
| `DownloadAttachment` | Download an attachment file |
| **Checklists** | |
| `ListChecklistItems` | List checklist items on an issue |
| `CreateChecklistItem` | Add a checklist item |
| `EditChecklistItem` | Edit a checklist item |
| `DeleteChecklistItem` | Delete a checklist item |
| **Links** | |
| `GetLinks` | List links for an issue |
| `CreateLink` | Create a link between issues |
| `DeleteLink` | Delete a link |
| **Transitions** | |
| `GetTransitions` | List available status transitions |
| `ExecuteTransition` | Execute a status transition |
| **Worklogs** | |
| `ListWorklogs` | List worklogs on an issue |
| `CreateWorklog` | Add a worklog entry |
| `EditWorklog` | Edit a worklog entry |
| `DeleteWorklog` | Delete a worklog entry |
| **Changelog** | |
| `GetChangelog` | Get the change history of an issue |

## See Also

- [ExampleIssuesService_Get](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Get)
- [ExampleIssuesService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Create)
- [ExampleIssuesService_Search](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search)
- [ExampleIssuesService_Search (pagination)](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.Search-pagination)
- [ExampleIssuesService_ListComments](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.ListComments)
- [ExampleIssuesService_UploadAttachment](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-IssuesService.UploadAttachment)
- [Authentication](auth.md)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
