# Queues

The Queues service manages Yandex Tracker queues and their sub-resources: auto-actions, macros, triggers, permissions, versions, and tags.

## Usage

### Get a Queue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

queue, _, err := client.Queues.Get(context.Background(), "QUEUE", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*queue.Key, *queue.Name)
```

### Create a Queue

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

queue, _, err := client.Queues.Create(context.Background(), &tracker.QueueCreateRequest{
    Key:  tracker.Ptr("NEWQUEUE"),
    Name: tracker.Ptr("New Queue"),
    Lead: tracker.Ptr("admin"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*queue.Key)
```

### List Queues

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

queues, _, err := client.Queues.List(context.Background(), nil)
if err != nil {
    log.Fatal(err)
}

for _, q := range queues {
    fmt.Println(*q.Key, *q.Name)
}
```

### List Macros

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

macros, _, err := client.Queues.ListMacros(context.Background(), "QUEUE")
if err != nil {
    log.Fatal(err)
}

for _, m := range macros {
    fmt.Println(*m.Name)
}
```

## Methods

| Method | Description |
|--------|-------------|
| **Core** | |
| `Create` | Create a new queue |
| `Get` | Get a queue by key |
| `List` | List all queues |
| `Delete` | Delete a queue |
| `Restore` | Restore a deleted queue |
| **Auto-actions** | |
| `ListAutoActions` | List auto-actions in a queue |
| `GetAutoAction` | Get an auto-action by ID |
| `CreateAutoAction` | Create an auto-action |
| `UpdateAutoAction` | Update an auto-action |
| **Macros** | |
| `ListMacros` | List macros in a queue |
| `GetMacro` | Get a macro by ID |
| `CreateMacro` | Create a macro |
| `EditMacro` | Edit a macro |
| `DeleteMacro` | Delete a macro |
| **Permissions** | |
| `GetPermissions` | Get queue permissions |
| `UpdatePermissions` | Update queue permissions |
| `CheckPermission` | Check a specific permission |
| **Triggers** | |
| `ListTriggers` | List triggers in a queue |
| `GetTrigger` | Get a trigger by ID |
| `CreateTrigger` | Create a trigger |
| `UpdateTrigger` | Update a trigger |
| **Versions** | |
| `ListVersions` | List versions in a queue |
| **Tags** | |
| `ListTags` | List tags in a queue |

## See Also

- [ExampleQueuesService_Get](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-QueuesService.Get)
- [ExampleQueuesService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-QueuesService.Create)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
