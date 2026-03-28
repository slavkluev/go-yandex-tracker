# Fields

The Fields service manages global and local issue fields. Global fields apply across all queues, while local fields are scoped to a specific queue.

## Usage

### List Global Fields

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

fields, _, err := client.Fields.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, field := range fields {
    fmt.Println(*field.ID, *field.Name)
}
```

### Create a Global Field

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

field, _, err := client.Fields.Create(context.Background(), &tracker.FieldCreateRequest{
    ID: tracker.Ptr("myField"),
    Name: &tracker.FieldName{
        EN: tracker.Ptr("My Field"),
        RU: tracker.Ptr("Мое поле"),
    },
    Category: tracker.Ptr("000000000000000000000002"),
    Type:     tracker.Ptr("ru.yandex.startrek.core.fields.StringFieldType"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*field.ID)
```

### List Local Fields

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

fields, _, err := client.Fields.ListLocal(context.Background(), "QUEUE")
if err != nil {
    log.Fatal(err)
}

for _, field := range fields {
    fmt.Println(*field.ID, *field.Name)
}
```

## Methods

**Global Fields**

| Method | Description |
|--------|-------------|
| `List` | List all global fields |
| `Get` | Get a global field by ID |
| `Create` | Create a global field |
| `Edit` | Edit a global field (requires version for optimistic locking) |

**Local Fields**

| Method | Description |
|--------|-------------|
| `ListLocal` | List local fields in a queue |
| `GetLocal` | Get a local field by key |
| `CreateLocal` | Create a local field in a queue |
| `EditLocal` | Edit a local field |

## See Also

- [ExampleFieldsService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-FieldsService.List)
- [ExampleFieldsService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-FieldsService.Create)
- [Error Handling](errors.md)
