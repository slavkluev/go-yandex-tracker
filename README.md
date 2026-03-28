# go-yandex-tracker

[![CI](https://github.com/slavkluev/go-yandex-tracker/actions/workflows/ci.yml/badge.svg)](https://github.com/slavkluev/go-yandex-tracker/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/slavkluev/go-yandex-tracker.svg)](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go client library for the [Yandex Tracker API](https://yandex.ru/support/tracker/en/).

## Features

- Zero dependencies -- only Go standard library
- 17 services covering 109+ API endpoints
- Idiomatic Go: `context.Context` on every call, typed errors, pointer fields with `Ptr[T]()` helper
- Page-based and scroll-token pagination
- Two auth modes: OAuth token and IAM token

## Installation

```bash
go get github.com/slavkluev/go-yandex-tracker
```

## Quick Start

### Authentication

**OAuth token** (for Yandex ID users):

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)
```

**IAM token** (for Yandex Cloud users):

```go
client := tracker.NewClient(
    tracker.WithIAMToken("your-iam-token"),
    tracker.WithCloudOrgID("your-cloud-org-id"),
)
```

### Get an Issue

```go
issue, _, err := client.Issues.Get(ctx, "QUEUE-1", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(*issue.Key, *issue.Summary)
```

### Create an Issue

```go
issue, _, err := client.Issues.Create(ctx, &tracker.IssueRequest{
    Summary: tracker.Ptr("Fix login page"),
    Queue:   tracker.Ptr("QUEUE"),
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(*issue.Key)
```

## Services

### Core

| Service | Description | Methods |
|---------|-------------|---------|
| [Issues](docs/issues.md) | Create, read, update, search issues; comments, attachments, worklogs, links, transitions, checklists | 31 |
| [Queues](docs/queues.md) | Queue management; triggers, auto-actions, macros, access | 23 |
| [Fields](docs/fields.md) | Issue field categories and configurations | 8 |
| [Components](docs/components.md) | Queue component management | 5 |

### Projects and Planning

| Service | Description | Methods |
|---------|-------------|---------|
| [Projects](docs/projects.md) | Projects, portfolios, and goals (entities) | 25 |
| [Boards](docs/boards.md) | Agile boards and columns | 9 |
| [Sprints](docs/sprints.md) | Sprint management | 3 |
| [Filters](docs/filters.md) | Saved search filters | 5 |

### Import and Automation

| Service | Description | Methods |
|---------|-------------|---------|
| [Import](docs/import.md) | Bulk data import | 5 |
| [Bulk Change](docs/bulkchange.md) | Bulk issue operations | 4 |
| [External Links](docs/externallinks.md) | External application links | 4 |

### Users and Reference

| Service | Description | Methods |
|---------|-------------|---------|
| [Users](docs/users.md) | User accounts | 3 |
| [Issue Types](docs/issuetypes.md) | Issue type definitions | 1 |
| [Statuses](docs/statuses.md) | Workflow statuses | 1 |
| [Resolutions](docs/resolutions.md) | Issue resolutions | 1 |
| [Priorities](docs/priorities.md) | Issue priorities | 1 |

### Dashboards

| Service | Description | Methods |
|---------|-------------|---------|
| [Dashboards](docs/dashboards.md) | Dashboards and widgets | 2 |

## Documentation

- [Authentication](docs/auth.md) -- OAuth and IAM token setup
- [Error Handling](docs/errors.md) -- typed errors and status codes
- [Pagination](docs/pagination.md) -- page-based and scroll-token patterns

For full API reference, see the [Go documentation](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker).

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This library is distributed under the MIT license. See the [LICENSE](LICENSE) file.
