# Sprints

The Sprints service manages sprints on agile boards.

## Usage

### List Sprints

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

sprints, _, err := client.Sprints.List(context.Background(), 1)
if err != nil {
    log.Fatal(err)
}

for _, s := range sprints {
    fmt.Println(*s.Name)
}
```

### Create a Sprint

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

sprint, _, err := client.Sprints.Create(context.Background(), &tracker.SprintCreateRequest{
    Name:  tracker.Ptr("Sprint 1"),
    Board: &tracker.BoardRef{ID: tracker.Ptr("1")},
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*sprint.Name)
```

## Methods

| Method | Description |
|--------|-------------|
| `List` | List sprints on a board |
| `Create` | Create a sprint |
| `Get` | Get a sprint by ID |

## See Also

- [ExampleSprintsService_List](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-SprintsService.List)
- [ExampleSprintsService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-SprintsService.Create)
- [Error Handling](errors.md)
