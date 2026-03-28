# Dashboards

The Dashboards service creates Yandex Tracker dashboards and adds widgets to them.

## Usage

### Create a Dashboard

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

dashboard, _, err := client.Dashboards.Create(context.Background(), &tracker.DashboardCreateRequest{
    Name: tracker.Ptr("My Dashboard"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*dashboard.Name)
```

### Create a Widget

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

widget, _, err := client.Dashboards.CreateWidget(context.Background(), 1, "cycleTime", &tracker.WidgetCreateRequest{})
if err != nil {
    log.Fatal(err)
}

fmt.Println(*widget.ID)
```

## Methods

| Method | Description |
|--------|-------------|
| `Create` | Create a dashboard |
| `CreateWidget` | Add a widget to a dashboard |

## See Also

- [ExampleDashboardsService_Create](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-DashboardsService.Create)
- [Error Handling](errors.md)
