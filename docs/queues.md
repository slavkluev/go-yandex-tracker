# Queues

The Queues service manages Yandex Tracker queues and their sub-resources: fields, workflows, components, auto-actions, macros, triggers, permissions, versions, and tags.

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

### Get a Queue with Expanded Data

`QueueGetOptions.Expand` asks Tracker for more data about the queue. Each value fills one `Queue` field, and `all` fills them all:

| `expand` | `Queue` field |
|----------|---------------|
| `team` | `TeamUsers` |
| `types` | `IssueTypes` |
| `versions` | `Versions` |
| `components` | `Components` |
| `workflows` | `Workflows` (workflow ID to its issue types) |
| `fields` | `Fields` |
| `issueTypesConfig` | `IssueTypesConfig` (issue type, workflow, resolutions) |

The expanded fields hold references (`Self`, `ID`, `Display`, and `Key` where the resource has one), not full resources. Tracker also accepts `projects`, but documents no shape for it, so `Queue` does not decode it.

```go
queue, _, err := client.Queues.Get(context.Background(), "QUEUE", &tracker.QueueGetOptions{
    Expand: "issueTypesConfig",
})
if err != nil {
    log.Fatal(err)
}

for _, c := range queue.IssueTypesConfig {
    fmt.Println(*c.IssueType.Key, *c.Workflow.ID)
}
```

### List Workflows and Their Transitions

`ListWorkflows` maps each workflow ID of a queue to its issue types. `Workflows.Get` returns a workflow's steps and transitions.

```go
workflows, _, err := client.Queues.ListWorkflows(context.Background(), "QUEUE")
if err != nil {
    log.Fatal(err)
}

for id := range workflows {
    workflow, _, err := client.Workflows.Get(context.Background(), id)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(id, len(workflow.Steps))
}
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
| `Get` | Get a queue by key; `expand` adds team, types, versions, components, workflows, fields, and issue type config |
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
| **Fields** | |
| `ListFields` | List the fields a queue configures, with per-queue `required` and options (often a subset, can be empty) |
| **Workflows** | |
| `ListWorkflows` | List a queue's workflows, each with its issue types |
| **Components** | |
| `ListComponents` | List components in a queue; `QueueComponentsListOptions.Fields` limits the returned fields |
| **Permissions** | |
| `UpdatePermissions` | Update queue permissions |
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
- [Workflows](workflows.md)
- [Components](components.md)
- [Fields](fields.md) -- global and local fields
- [Error Handling](errors.md)
- [Pagination](pagination.md)
