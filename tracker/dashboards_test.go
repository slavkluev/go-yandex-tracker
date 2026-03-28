package tracker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestDashboardsService_Create(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"name":"My Dashboard","layout":"two-columns"}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/dashboards/123",
			"id": 123,
			"version": 1,
			"name": "My Dashboard",
			"layout": "two-columns",
			"owner": {
				"self": "https://api.tracker.yandex.net/v3/users/1",
				"id": "user1",
				"display": "User One"
			},
			"createdBy": {
				"self": "https://api.tracker.yandex.net/v3/users/1",
				"id": "user1",
				"display": "User One"
			},
			"createdAt": "2026-03-26T12:00:00.000+0000"
		}`)
	})

	ctx := context.Background()
	dashboard, resp, err := client.Dashboards.Create(ctx, &DashboardCreateRequest{
		Name:   Ptr("My Dashboard"),
		Layout: Ptr("two-columns"),
	})
	if err != nil {
		t.Fatalf("Dashboards.Create returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	want := &Dashboard{
		Self:    Ptr("https://api.tracker.yandex.net/v3/dashboards/123"),
		ID:      Ptr(123),
		Version: Ptr(1),
		Name:    Ptr("My Dashboard"),
		Layout:  Ptr("two-columns"),
		Owner: &User{
			Self:    Ptr("https://api.tracker.yandex.net/v3/users/1"),
			ID:      Ptr("user1"),
			Display: Ptr("User One"),
		},
		CreatedBy: &User{
			Self:    Ptr("https://api.tracker.yandex.net/v3/users/1"),
			ID:      Ptr("user1"),
			Display: Ptr("User One"),
		},
		CreatedAt: newTestTimestamp(t, "2026-03-26T12:00:00.000+0000"),
	}

	if !reflect.DeepEqual(dashboard, want) {
		t.Errorf("Dashboards.Create returned %+v, want %+v", dashboard, want)
	}
}

func TestDashboardsService_CreateWidget(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/dashboards/{dashboardID}/widgets/{widgetType}", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://api.tracker.yandex.net/v3/widgets/456",
			"id": 456,
			"version": 1,
			"description": "Cycle Time Widget",
			"createdBy": {
				"self": "https://api.tracker.yandex.net/v3/users/1",
				"id": "user1",
				"display": "User One"
			},
			"dashboard": {
				"id": "123",
				"display": "My Dashboard"
			},
			"datasetInfo": {
				"status": "building",
				"buildStartedAt": "2026-03-26T12:00:00.000+0000"
			},
			"mode": "common-lines",
			"autoUpdatable": true
		}`)
	})

	ctx := context.Background()
	widget, resp, err := client.Dashboards.CreateWidget(ctx, 123, "cycleTime", &WidgetCreateRequest{
		Description: Ptr("Cycle Time Widget"),
		Parameters: map[string]any{
			"mode":          "common-lines",
			"autoUpdatable": true,
		},
	})
	if err != nil {
		t.Fatalf("Dashboards.CreateWidget returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if got, want := *widget.Self, "https://api.tracker.yandex.net/v3/widgets/456"; got != want {
		t.Errorf("Widget.Self = %q, want %q", got, want)
	}
	if got, want := *widget.ID, 456; got != want {
		t.Errorf("Widget.ID = %d, want %d", got, want)
	}
	if got, want := *widget.Version, 1; got != want {
		t.Errorf("Widget.Version = %d, want %d", got, want)
	}
	if got, want := *widget.Description, "Cycle Time Widget"; got != want {
		t.Errorf("Widget.Description = %q, want %q", got, want)
	}
	if got, want := *widget.CreatedBy.ID, "user1"; got != want {
		t.Errorf("Widget.CreatedBy.ID = %q, want %q", got, want)
	}
	if got, want := *widget.Dashboard.ID, "123"; got != want {
		t.Errorf("Widget.Dashboard.ID = %q, want %q", got, want)
	}
	if got, want := *widget.Dashboard.Display, "My Dashboard"; got != want {
		t.Errorf("Widget.Dashboard.Display = %q, want %q", got, want)
	}
	if got, want := *widget.DatasetInfo.Status, "building"; got != want {
		t.Errorf("Widget.DatasetInfo.Status = %q, want %q", got, want)
	}
	if got, want := *widget.DatasetInfo.BuildStartedAt, "2026-03-26T12:00:00.000+0000"; got != want {
		t.Errorf("Widget.DatasetInfo.BuildStartedAt = %q, want %q", got, want)
	}

	// Verify type-specific parameters are captured in Parameters map.
	if got, want := widget.Parameters["mode"], "common-lines"; got != want {
		t.Errorf("Widget.Parameters[mode] = %v, want %v", got, want)
	}
	if got, want := widget.Parameters["autoUpdatable"], true; got != want {
		t.Errorf("Widget.Parameters[autoUpdatable] = %v, want %v", got, want)
	}
}

// newTestTimestamp parses a Yandex Tracker timestamp string and returns
// a *Timestamp for use in test expected values. This ensures the timezone
// representation matches what json.Unmarshal produces.
func newTestTimestamp(t *testing.T, s string) *Timestamp {
	t.Helper()
	parsed, err := time.Parse("2006-01-02T15:04:05.000-0700", s)
	if err != nil {
		t.Fatalf("newTestTimestamp: failed to parse %q: %v", s, err)
	}
	return &Timestamp{Time: parsed}
}
