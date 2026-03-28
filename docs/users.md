# Users

The Users service provides information about Yandex Tracker users: the authenticated user, individual users by ID, and user listing.

## Usage

### Get Current User

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

user, _, err := client.Users.Myself(context.Background())
if err != nil {
    log.Fatal(err)
}

fmt.Println(*user.Display)
```

### Get a User by ID

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

user, _, err := client.Users.Get(context.Background(), "user123")
if err != nil {
    log.Fatal(err)
}

fmt.Println(*user.Display)
```

## Methods

| Method | Description |
|--------|-------------|
| `Myself` | Get the authenticated user |
| `Get` | Get a user by ID or login |
| `List` | List users in the organization |

## See Also

- [ExampleUsersService_Myself](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-UsersService.Myself)
- [ExampleUsersService_Get](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-UsersService.Get)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
