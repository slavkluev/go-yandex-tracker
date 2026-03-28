# Boards

The Boards service manages Yandex Tracker agile boards and their columns. Board mutations use optimistic locking via a version parameter.

## Usage

### List Boards

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

boards, _, err := client.Boards.List(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, b := range boards {
    fmt.Println(*b.Name)
}
```

### Edit a Board

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

// The version parameter is used for optimistic locking via the If-Match header.
// Get the current version from the board's Version field first.
board, _, err := client.Boards.Edit(context.Background(), 1, 5, &tracker.BoardEditRequest{
    Name: tracker.Ptr("Updated Board Name"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*board.Name)
```

### Create a Column

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

// Version for optimistic locking -- get from board's Version field.
column, _, err := client.Boards.CreateColumn(context.Background(), 1, 5, &tracker.ColumnCreateRequest{
    Name: tracker.Ptr("In Progress"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*column.Name)
```

## Methods

| Method | Description |
|--------|-------------|
| **Core** | |
| `Create` | Create a new board |
| `Get` | Get a board by ID |
| `List` | List all boards |
| `Edit` | Edit a board (requires version for optimistic locking) |
| `Delete` | Delete a board |
| **Columns** | |
| `ListColumns` | List columns on a board |
| `CreateColumn` | Create a column (requires version) |
| `EditColumn` | Edit a column (requires version) |
| `DeleteColumn` | Delete a column (requires version) |

## See Also

- [ExampleBoardsService_Edit](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-BoardsService.Edit)
- [ExampleBoardsService_CreateColumn](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-BoardsService.CreateColumn)
- [Error Handling](errors.md)
