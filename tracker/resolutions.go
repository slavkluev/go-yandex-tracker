package tracker

import "context"

// List returns all resolutions.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-resolutions
func (s *ResolutionsService) List(ctx context.Context) ([]*Resolution, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/resolutions", nil)
	if err != nil {
		return nil, nil, err
	}

	var resolutions []*Resolution
	resp, err := s.client.Do(ctx, req, &resolutions)
	if err != nil {
		return nil, resp, err
	}

	return resolutions, resp, nil
}
