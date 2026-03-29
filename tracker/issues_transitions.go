package tracker

import (
	"context"
	"fmt"
)

// GetTransitions returns the available transitions for an issue.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/get-transitions
func (s *IssuesService) GetTransitions(ctx context.Context, issueKey string) ([]*Transition, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/transitions", issueKey)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var transitions []*Transition
	resp, err := s.client.Do(ctx, req, &transitions)
	if err != nil {
		return nil, resp, err
	}

	return transitions, resp, nil
}

// ExecuteTransition executes a transition on an issue.
// The transition parameter specifies optional fields to set during the transition.
// If transition is nil, the transition is executed without additional parameters.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/issues/new-transition
func (s *IssuesService) ExecuteTransition(ctx context.Context, issueKey, transitionID string, transition *TransitionRequest) ([]*Transition, *Response, error) {
	u := fmt.Sprintf("v2/issues/%v/transitions/%v/_execute", issueKey, transitionID)

	req, err := s.client.NewRequest("POST", u, transition)
	if err != nil {
		return nil, nil, err
	}

	var transitions []*Transition
	resp, err := s.client.Do(ctx, req, &transitions)
	if err != nil {
		return nil, resp, err
	}

	return transitions, resp, nil
}
