// Package tracker provides a Go client for the Yandex Tracker API.
//
// Usage:
//
//	client := tracker.NewClient(
//		tracker.WithOAuthToken("your-oauth-token"),
//		tracker.WithOrgID("your-org-id"),
//	)
//
//	// Get an issue
//	issue, _, err := client.Issues.Get(ctx, "QUEUE-1", nil)
//
//	// Create an issue
//	issue, _, err := client.Issues.Create(ctx, &tracker.IssueRequest{
//		Summary: tracker.Ptr("New issue"),
//		Queue:   tracker.Ptr("QUEUE"),
//	})
//
//	// Search issues
//	issues, _, err := client.Issues.Search(ctx, &tracker.IssueSearchRequest{
//		Filter: map[string]any{"queue": "QUEUE"},
//	}, nil)
//
// Authentication:
//
// The client supports two authentication methods:
//   - OAuth token: tracker.WithOAuthToken("token") sets the Authorization: OAuth header
//   - IAM token: tracker.WithIAMToken("token") sets the Authorization: Bearer header
//
// Organization:
//
// Yandex Tracker requires an organization identifier:
//   - Yandex 360: tracker.WithOrgID("org-id") sets the X-Org-Id header
//   - Yandex Cloud: tracker.WithCloudOrgID("cloud-org-id") sets the X-Cloud-Org-Id header
//
// The Yandex Tracker API documentation is available at
// https://yandex.ru/support/tracker/en/
package tracker
