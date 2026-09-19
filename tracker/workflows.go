package tracker

import (
	"context"
	"fmt"
)

// Get fetches a workflow by its ID, with its steps (statuses) and the actions
// (transitions) available from each step. QueuesService.ListWorkflows returns
// the IDs of a queue's workflows.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api/queues/workflows/get-workflow
func (s *WorkflowsService) Get(ctx context.Context, workflowID string) (*Workflow, *Response, error) {
	u := fmt.Sprintf("v3/workflows/%v", workflowID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	workflow := new(Workflow)
	resp, err := s.client.Do(ctx, req, workflow)
	if err != nil {
		return nil, resp, err
	}

	return workflow, resp, nil
}
