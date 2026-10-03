package tracker

import (
	"context"
	"fmt"
	"iter"
)

// Myself returns information about the currently authenticated user.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/users/get-user-info
func (s *UsersService) Myself(ctx context.Context) (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "v3/myself", nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// Get fetches a user by their login or user ID.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/users/get-user
func (s *UsersService) Get(ctx context.Context, userID string) (*User, *Response, error) {
	u := fmt.Sprintf("v3/users/%v", userID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// List returns a paginated list of users in the organization.
// Pass nil for opts to use default pagination.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/users/get-users
func (s *UsersService) List(ctx context.Context, opts *UserListOptions) ([]*User, *Response, error) {
	u := "v3/users"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// ListIter yields the users of every page, starting with the page List
// returns for opts, which it copies when called. On an error it yields the
// error once and stops.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/users/get-users
func (s *UsersService) ListIter(ctx context.Context, opts *UserListOptions) iter.Seq2[*User, error] {
	return pageIter(ctx, opts, s.List)
}
