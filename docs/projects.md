# Projects

The Projects service manages Yandex Tracker entities -- projects, portfolios, and goals. The underlying Go type is `EntitiesService`, accessed via `client.Entities`. Methods take an `EntityType` parameter (`tracker.EntityTypeProject`, `tracker.EntityTypePortfolio`, or `tracker.EntityTypeGoal`) to specify the entity kind. This document focuses on project usage.

## Usage

### Create a Project

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

entity, _, err := client.Entities.Create(context.Background(), tracker.EntityTypeProject, &tracker.EntityCreateRequest{
    Fields: &tracker.EntityCreateFields{
        Summary:    tracker.Ptr("Q1 Release"),
        TeamAccess: tracker.Ptr(true),
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*entity.ID)
```

### Search Projects

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

result, _, err := client.Entities.Search(context.Background(), tracker.EntityTypeProject, &tracker.EntitySearchRequest{
    Input: tracker.Ptr("project"),
}, nil)
if err != nil {
    log.Fatal(err)
}

for _, e := range result.Values {
    fmt.Println(*e.ID)
}
```

### Get a Project

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

entity, _, err := client.Entities.Get(context.Background(), tracker.EntityTypeProject, "6", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*entity.ID)
```

### Update a Project

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

entity, _, err := client.Entities.Update(context.Background(), tracker.EntityTypeProject, "6", &tracker.EntityUpdateRequest{
    Fields: &tracker.EntityCreateFields{
        Summary: tracker.Ptr("Updated Project Name"),
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*entity.ID)
```

## Methods

| Method | Description |
|--------|-------------|
| **Core** | |
| `Create` | Create a new entity (project, portfolio, or goal) |
| `Get` | Get an entity by ID |
| `Update` | Update an entity |
| `Delete` | Delete an entity |
| `Search` | Search entities with filters |
| `BulkChange` | Bulk update entity fields |
| `GetEvents` | Get entity event history |
| **Access** | |
| `GetAccess` | Get entity access permissions |
| `UpdateAccess` | Update entity access permissions |
| **Attachments** | |
| `ListAttachments` | List attachments on an entity |
| `AttachFile` | Attach a previously uploaded file |
| `DeleteAttachment` | Delete an attachment |
| **Checklists** | |
| `CreateChecklistItem` | Add a checklist item |
| `EditChecklistItem` | Edit a checklist item |
| `MoveChecklistItem` | Move a checklist item |
| `DeleteChecklistItem` | Delete a checklist item |
| `DeleteChecklist` | Delete entire checklist |
| **Comments** | |
| `ListComments` | List comments on an entity |
| `CreateComment` | Add a comment |
| `GetComment` | Get a comment by ID |
| `EditComment` | Edit a comment |
| `DeleteComment` | Delete a comment |
| **Links** | |
| `ListLinks` | List entity links |
| `CreateLink` | Create entity links |
| `DeleteLink` | Delete an entity link |

## See Also

- [ExampleEntitiesService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-EntitiesService.Create)
- [ExampleEntitiesService_Search](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-EntitiesService.Search)
- [Error Handling](errors.md)
