package tracker

import (
	"context"
	"fmt"
)

// ListWorkflows returns the workflows of a queue, as a map from each
// workflow ID to the issue types that follow that workflow. Pass a workflow
// ID to WorkflowsService.Get for its steps and transitions.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api/queues/workflows/get-queue-workflows
func (s *QueuesService) ListWorkflows(ctx context.Context, queueKey string) (map[string][]*IssueType, *Response, error) {
	u := fmt.Sprintf("v3/queues/%v/workflows", queueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var workflows map[string][]*IssueType
	resp, err := s.client.Do(ctx, req, &workflows)
	if err != nil {
		return nil, resp, err
	}

	return workflows, resp, nil
}
