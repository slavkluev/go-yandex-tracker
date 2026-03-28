# Bulk Change

The Bulk Change service performs asynchronous bulk operations on issues: move, update fields, and execute transitions.

## Usage

### Bulk Move Issues

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

// Initiate bulk move -- returns immediately with operation ID.
op, _, err := client.BulkChange.Move(context.Background(), &tracker.BulkMoveRequest{
    Queue:  tracker.Ptr("TARGET"),
    Issues: []string{"SOURCE-1", "SOURCE-2", "SOURCE-3"},
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*op.ID)
```

### Check Operation Status

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

// Poll for completion using the operation ID.
status, _, err := client.BulkChange.GetStatus(context.Background(), *op.ID)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*status.Status)
```

## Methods

| Method | Description |
|--------|-------------|
| `Move` | Bulk move issues to another queue |
| `Update` | Bulk update issue fields |
| `Transition` | Bulk execute status transitions |
| `GetStatus` | Get the status of a bulk operation |

## See Also

- [ExampleBulkChangeService_Move](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-BulkChangeService.Move)
- [Error Handling](errors.md)
