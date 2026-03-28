package tracker

import (
	"net/url"
	"testing"
)

func TestAddOptions(t *testing.T) {
	t.Run("nil opts", func(t *testing.T) {
		got, err := addOptions("issues", nil)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}
		if got != "issues" {
			t.Errorf("addOptions = %q, want %q", got, "issues")
		}
	})

	t.Run("nil pointer to struct", func(t *testing.T) {
		var opts *ListOptions
		got, err := addOptions("issues", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}
		if got != "issues" {
			t.Errorf("addOptions = %q, want %q", got, "issues")
		}
	})

	t.Run("empty struct omits zero values", func(t *testing.T) {
		opts := &ListOptions{}
		got, err := addOptions("issues", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}
		if got != "issues" {
			t.Errorf("addOptions = %q, want %q", got, "issues")
		}
	})

	t.Run("ListOptions with values", func(t *testing.T) {
		opts := &ListOptions{Page: 2, PerPage: 50}
		got, err := addOptions("issues", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("page"); v != "2" {
			t.Errorf("page = %q, want %q", v, "2")
		}
		if v := u.Query().Get("perPage"); v != "50" {
			t.Errorf("perPage = %q, want %q", v, "50")
		}
	})

	t.Run("IssueSearchOptions with embedded ListOptions", func(t *testing.T) {
		opts := &IssueSearchOptions{
			ListOptions: ListOptions{Page: 1, PerPage: 10},
			Expand:      "transitions",
		}
		got, err := addOptions("issues/_search", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("page"); v != "1" {
			t.Errorf("page = %q, want %q", v, "1")
		}
		if v := u.Query().Get("perPage"); v != "10" {
			t.Errorf("perPage = %q, want %q", v, "10")
		}
		if v := u.Query().Get("expand"); v != "transitions" {
			t.Errorf("expand = %q, want %q", v, "transitions")
		}
	})

	t.Run("ScrollSearchOptions", func(t *testing.T) {
		opts := &ScrollSearchOptions{
			ScrollType:  "sorted",
			PerScroll:   100,
			ScrollTTLMs: 60000,
		}
		got, err := addOptions("issues/_search", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("scrollType"); v != "sorted" {
			t.Errorf("scrollType = %q, want %q", v, "sorted")
		}
		if v := u.Query().Get("perScroll"); v != "100" {
			t.Errorf("perScroll = %q, want %q", v, "100")
		}
		if v := u.Query().Get("scrollTTLMillis"); v != "60000" {
			t.Errorf("scrollTTLMillis = %q, want %q", v, "60000")
		}
	})

	t.Run("ScrollSearchOptions omits zero values", func(t *testing.T) {
		opts := &ScrollSearchOptions{
			ScrollType: "sorted",
		}
		got, err := addOptions("issues/_search", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("scrollType"); v != "sorted" {
			t.Errorf("scrollType = %q, want %q", v, "sorted")
		}
		if v := u.Query().Get("perScroll"); v != "" {
			t.Errorf("perScroll should be omitted, got %q", v)
		}
		if v := u.Query().Get("scrollTTLMillis"); v != "" {
			t.Errorf("scrollTTLMillis should be omitted, got %q", v)
		}
	})

	t.Run("IssueSearchOptions with only embedded fields", func(t *testing.T) {
		opts := &IssueSearchOptions{
			ListOptions: ListOptions{Page: 3, PerPage: 25},
		}
		got, err := addOptions("issues/_search", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("page"); v != "3" {
			t.Errorf("page = %q, want %q", v, "3")
		}
		if v := u.Query().Get("perPage"); v != "25" {
			t.Errorf("perPage = %q, want %q", v, "25")
		}
		// expand should not be present
		if v := u.Query().Get("expand"); v != "" {
			t.Errorf("expand should be omitted, got %q", v)
		}
	})

	t.Run("ChangelogOptions", func(t *testing.T) {
		opts := &ChangelogOptions{
			ID:      "abc123",
			PerPage: 20,
		}
		got, err := addOptions("issues/KEY-1/changelog", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("id"); v != "abc123" {
			t.Errorf("id = %q, want %q", v, "abc123")
		}
		if v := u.Query().Get("perPage"); v != "20" {
			t.Errorf("perPage = %q, want %q", v, "20")
		}
	})

	t.Run("IssueMoveOptions with bool pointers", func(t *testing.T) {
		opts := &IssueMoveOptions{
			Queue:         "NEWQUEUE",
			MoveAllFields: Ptr(true),
			InitialStatus: Ptr(false),
		}
		got, err := addOptions("issues/KEY-1/_move", opts)
		if err != nil {
			t.Fatalf("addOptions returned error: %v", err)
		}

		u, err := url.Parse(got)
		if err != nil {
			t.Fatalf("failed to parse result URL: %v", err)
		}

		if v := u.Query().Get("queue"); v != "NEWQUEUE" {
			t.Errorf("queue = %q, want %q", v, "NEWQUEUE")
		}
		if v := u.Query().Get("moveAllFields"); v != "true" {
			t.Errorf("moveAllFields = %q, want %q", v, "true")
		}
		if v := u.Query().Get("initialStatus"); v != "false" {
			t.Errorf("initialStatus = %q, want %q", v, "false")
		}
		// expandTransitions is nil, should be omitted
		if v := u.Query().Get("expandTransitions"); v != "" {
			t.Errorf("expandTransitions should be omitted, got %q", v)
		}
	})
}
