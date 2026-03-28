package tracker

import (
	"context"
	"fmt"
)

// Create creates a new dashboard.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/dashboards/create-dashboard
func (s *DashboardsService) Create(ctx context.Context, dashboard *DashboardCreateRequest) (*Dashboard, *Response, error) {
	req, err := s.client.NewRequest("POST", "v3/dashboards/", dashboard)
	if err != nil {
		return nil, nil, err
	}

	d := new(Dashboard)
	resp, err := s.client.Do(ctx, req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// CreateWidget creates a widget of the specified type on a dashboard.
// The widgetType parameter is used to construct the URL path (e.g., "cycleTime").
// The request body contains widget-type-specific parameters.
//
// Currently the only documented widget type is "cycleTime".
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/api-ref/dashboards/create-widget
func (s *DashboardsService) CreateWidget(ctx context.Context, dashboardID int, widgetType string, widget *WidgetCreateRequest) (*Widget, *Response, error) {
	u := fmt.Sprintf("v3/dashboards/%v/widgets/%v", dashboardID, widgetType)

	req, err := s.client.NewRequest("POST", u, widget)
	if err != nil {
		return nil, nil, err
	}

	w := new(Widget)
	resp, err := s.client.Do(ctx, req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
}
