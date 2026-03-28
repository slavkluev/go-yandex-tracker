package tracker

import "context"

// List returns all issue types.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/admin/get-issue-types
func (s *IssueTypesService) List(ctx context.Context) ([]*IssueType, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/issuetypes", nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(ctx, req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}
