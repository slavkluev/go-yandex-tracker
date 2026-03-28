package tracker

import (
	"context"
	"fmt"
)

// ListLinks returns the links for an entity.
// Unlike issue links which return []*IssueLink, entity links return
// []*EntityLinkResponse containing the link type and field values
// of the linked entity.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/links/get-links
func (s *EntitiesService) ListLinks(ctx context.Context, entityType EntityType, id string) ([]*EntityLinkResponse, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/links", entityType, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var links []*EntityLinkResponse
	resp, err := s.client.Do(ctx, req, &links)
	if err != nil {
		return nil, resp, err
	}

	return links, resp, nil
}

// CreateLink creates links between entities.
// The links parameter is a slice of EntityLink, each specifying a
// relationship type and target entity ID. The API accepts an array
// body (not a single object), allowing multiple links to be created
// in one request.
//
// Returns the updated Entity reflecting the new links.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/links/add-links
func (s *EntitiesService) CreateLink(ctx context.Context, entityType EntityType, id string, links []*EntityLink) (*Entity, *Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/links", entityType, id)

	req, err := s.client.NewRequest("POST", u, links)
	if err != nil {
		return nil, nil, err
	}

	entity := new(Entity)
	resp, err := s.client.Do(ctx, req, entity)
	if err != nil {
		return nil, resp, err
	}

	return entity, resp, nil
}

// DeleteLink removes a link from an entity.
// Unlike issue link deletion which uses a link ID in the URL path,
// entity link deletion uses a query parameter "right" specifying the
// ID of the linked entity to unlink.
//
// The API returns HTTP 200 OK (not 204 No Content).
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/ru/api-ref/entities/links/delete-link
func (s *EntitiesService) DeleteLink(ctx context.Context, entityType EntityType, id string, opts *EntityLinkDeleteOptions) (*Response, error) {
	u := fmt.Sprintf("v3/entities/%v/%v/links", entityType, id)

	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
