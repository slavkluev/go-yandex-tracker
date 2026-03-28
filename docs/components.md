# Components

The Components service manages queue components for categorizing issues.

## Usage

### List Components

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

components, _, err := client.Components.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, c := range components {
    fmt.Println(*c.Name)
}
```

### Create a Component

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

component, _, err := client.Components.Create(context.Background(), &tracker.ComponentRequest{
    Name:  tracker.Ptr("Backend"),
    Queue: tracker.Ptr("QUEUE"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*component.Name)
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List all components |
| `Create` | Create a component |
| `Get` | Get a component by ID |
| `Edit` | Edit a component |
| `Delete` | Delete a component |

## See Also

- [ExampleComponentsService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ComponentsService.List)
- [ExampleComponentsService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-ComponentsService.Create)
- [Error Handling](errors.md)
