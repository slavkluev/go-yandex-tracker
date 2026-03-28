# Authentication

The Yandex Tracker API requires an authentication token and an organization ID. Two authentication modes are supported: OAuth tokens (for Yandex 360 organizations) and IAM tokens (for Yandex Cloud organizations).

## OAuth Token (Yandex 360)

Use `WithOAuthToken` and `WithOrgID` for Yandex 360 organizations. The token is sent as `Authorization: OAuth <token>` and the org ID as `X-Org-Id`.

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("y0_AgAAAAD..."),
    tracker.WithOrgID("12345"),
)
```

## IAM Token (Yandex Cloud)

Use `WithIAMToken` and `WithCloudOrgID` for Yandex Cloud organizations. The token is sent as `Authorization: Bearer <token>` and the org ID as `X-Cloud-Org-Id`.

```go
client := tracker.NewClient(
    tracker.WithIAMToken("t1.9euelZqM..."),
    tracker.WithCloudOrgID("bpf1234abc"),
)
```

## Custom HTTP Client

Pass a custom `http.Client` for timeouts, proxies, or transport-level middleware:

```go
httpClient := &http.Client{Timeout: 30 * time.Second}
client := tracker.NewClient(
    tracker.WithOAuthToken("your-token"),
    tracker.WithOrgID("your-org"),
    tracker.WithHTTPClient(httpClient),
)
```

## See Also

- [ExampleNewClient_oAuth](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-NewClient-oAuth)
- [ExampleNewClient_iAMToken](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-NewClient-iAMToken)
- [Error Handling](errors.md)
- [Pagination](pagination.md)
