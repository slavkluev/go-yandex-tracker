package tracker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://api.tracker.yandex.net/"
)

type authType int

const (
	authNone authType = iota
	authOAuth
	authIAM
)

// Client manages communication with the Yandex Tracker API.
type Client struct {
	client     *http.Client
	baseURL    *url.URL
	token      string
	authType   authType
	orgID      string
	cloudOrgID string

	common service // Reuse a single struct instead of allocating per service.

	// Issues handles communication with the issue related methods
	// of the Yandex Tracker API.
	Issues *IssuesService

	// Queues handles communication with the queue related methods
	// of the Yandex Tracker API.
	Queues *QueuesService

	// Fields handles communication with the issue field related methods
	// of the Yandex Tracker API.
	Fields *FieldsService

	// Components handles communication with the component related methods
	// of the Yandex Tracker API.
	Components *ComponentsService

	// IssueTypes handles communication with the issue type related methods
	// of the Yandex Tracker API.
	IssueTypes *IssueTypesService

	// Statuses handles communication with the status related methods
	// of the Yandex Tracker API.
	Statuses *StatusesService

	// Resolutions handles communication with the resolution related methods
	// of the Yandex Tracker API.
	Resolutions *ResolutionsService

	// Priorities handles communication with the priority related methods
	// of the Yandex Tracker API.
	Priorities *PrioritiesService

	// Users handles communication with the user related methods
	// of the Yandex Tracker API.
	Users *UsersService

	// BulkChange handles communication with the bulk change related methods
	// of the Yandex Tracker API.
	BulkChange *BulkChangeService

	// Boards handles communication with the board and column related methods
	// of the Yandex Tracker API.
	Boards *BoardsService

	// Sprints handles communication with the sprint related methods
	// of the Yandex Tracker API.
	Sprints *SprintsService

	// Entities handles communication with the entity related methods
	// (projects, portfolios, goals) of the Yandex Tracker API.
	Entities *EntitiesService

	// Filters handles communication with the saved filter related methods
	// of the Yandex Tracker API.
	Filters *FiltersService

	// ExternalLinks handles communication with the external application link related methods
	// of the Yandex Tracker API.
	ExternalLinks *ExternalLinksService

	// Import handles communication with the issue import related methods
	// of the Yandex Tracker API.
	Import *ImportService

	// Dashboards handles communication with the dashboard and widget
	// related methods of the Yandex Tracker API.
	Dashboards *DashboardsService
}

type service struct {
	client *Client
}

// IssuesService handles communication with the issue related methods
// of the Yandex Tracker API.
type IssuesService service

// QueuesService handles communication with the queue related methods
// of the Yandex Tracker API.
type QueuesService service

// FieldsService handles communication with the issue field related methods
// of the Yandex Tracker API.
type FieldsService service

// ComponentsService handles communication with the component related methods
// of the Yandex Tracker API.
type ComponentsService service

// IssueTypesService handles communication with the issue type related methods
// of the Yandex Tracker API.
type IssueTypesService service

// StatusesService handles communication with the status related methods
// of the Yandex Tracker API.
type StatusesService service

// ResolutionsService handles communication with the resolution related methods
// of the Yandex Tracker API.
type ResolutionsService service

// PrioritiesService handles communication with the priority related methods
// of the Yandex Tracker API.
type PrioritiesService service

// UsersService handles communication with the user related methods
// of the Yandex Tracker API.
type UsersService service

// BulkChangeService handles communication with the bulk change related methods
// of the Yandex Tracker API.
type BulkChangeService service

// BoardsService handles communication with the board and column related methods
// of the Yandex Tracker API.
type BoardsService service

// SprintsService handles communication with the sprint related methods
// of the Yandex Tracker API.
type SprintsService service

// EntitiesService handles communication with the entity related methods
// (projects, portfolios, goals) of the Yandex Tracker API.
type EntitiesService service

// FiltersService handles communication with the saved filter related methods
// of the Yandex Tracker API.
type FiltersService service

// ExternalLinksService handles communication with the external application link related methods
// of the Yandex Tracker API.
type ExternalLinksService service

// ImportService handles communication with the issue import related methods
// of the Yandex Tracker API.
//
// Import methods require admin permissions in the target queue.
// They preserve original creation timestamps (createdAt/createdBy) for historical data migration.
type ImportService service

// DashboardsService handles communication with the dashboard and widget
// related methods of the Yandex Tracker API.
type DashboardsService service

// Option configures a Client.
type Option func(*Client)

// WithOAuthToken sets the OAuth token for authentication.
// The token is sent as "Authorization: OAuth <token>" header.
func WithOAuthToken(token string) Option {
	return func(c *Client) {
		c.token = token
		c.authType = authOAuth
	}
}

// WithIAMToken sets the IAM token for authentication.
// The token is sent as "Authorization: Bearer <token>" header.
func WithIAMToken(token string) Option {
	return func(c *Client) {
		c.token = token
		c.authType = authIAM
	}
}

// WithOrgID sets the organization ID for Yandex 360.
// Sent as "X-Org-Id" header on every request.
func WithOrgID(orgID string) Option {
	return func(c *Client) {
		c.orgID = orgID
	}
}

// WithCloudOrgID sets the organization ID for Yandex Cloud.
// Sent as "X-Cloud-Org-Id" header on every request.
func WithCloudOrgID(cloudOrgID string) Option {
	return func(c *Client) {
		c.cloudOrgID = cloudOrgID
	}
}

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.client = httpClient
	}
}

// WithBaseURL sets the base URL for API requests.
// The URL must have a trailing slash.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		u, _ := url.Parse(baseURL)
		c.baseURL = u
	}
}

// NewClient creates a new Yandex Tracker API client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.baseURL == nil {
		c.baseURL, _ = url.Parse(defaultBaseURL)
	}
	c.initialize()
	return c
}

func (c *Client) initialize() {
	c.common.client = c
	c.Issues = (*IssuesService)(&c.common)
	c.Queues = (*QueuesService)(&c.common)
	c.Fields = (*FieldsService)(&c.common)
	c.Components = (*ComponentsService)(&c.common)
	c.IssueTypes = (*IssueTypesService)(&c.common)
	c.Statuses = (*StatusesService)(&c.common)
	c.Resolutions = (*ResolutionsService)(&c.common)
	c.Priorities = (*PrioritiesService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.BulkChange = (*BulkChangeService)(&c.common)
	c.Boards = (*BoardsService)(&c.common)
	c.Sprints = (*SprintsService)(&c.common)
	c.Entities = (*EntitiesService)(&c.common)
	c.Filters = (*FiltersService)(&c.common)
	c.ExternalLinks = (*ExternalLinksService)(&c.common)
	c.Import = (*ImportService)(&c.common)
	c.Dashboards = (*DashboardsService)(&c.common)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// which will be resolved against the baseURL of the Client. If body is
// specified, it is JSON encoded and included as the request body. If body
// is nil, no Content-Type header is set and no body is sent.
func (c *Client) NewRequest(method, urlStr string, body any) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL.String())
	}

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	//nolint:noctx // context is added later in Do/BareDo via req.WithContext(ctx); matches go-github pattern
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	switch c.authType {
	case authOAuth:
		req.Header.Set("Authorization", "OAuth "+c.token)
	case authIAM:
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	if c.orgID != "" {
		req.Header.Set("X-Org-Id", c.orgID)
	}
	if c.cloudOrgID != "" {
		req.Header.Set("X-Cloud-Org-Id", c.cloudOrgID)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error in case of an API error.
//
// If v is nil, the response body is not decoded.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	if err := checkResponse(resp); err != nil {
		return response, err
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // Empty body is not an error.
		}
		if decErr != nil {
			err = decErr
		}
	}

	return response, err
}

// BareDo sends an API request and returns the raw response. The response body
// is not closed; the caller is responsible for closing it.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	if err := checkResponse(resp); err != nil {
		resp.Body.Close()
		return response, err
	}

	return response, nil
}

// NewUploadRequest creates a multipart/form-data upload request.
// The file content is read from reader and sent with the given filename.
// Auth and org headers are set identically to NewRequest.
func (c *Client) NewUploadRequest(method, urlStr, filename string, reader io.Reader) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL.String())
	}

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, reader); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	//nolint:noctx // context is added later in Do/BareDo via req.WithContext(ctx); matches go-github pattern
	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	switch c.authType {
	case authOAuth:
		req.Header.Set("Authorization", "OAuth "+c.token)
	case authIAM:
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.orgID != "" {
		req.Header.Set("X-Org-Id", c.orgID)
	}
	if c.cloudOrgID != "" {
		req.Header.Set("X-Cloud-Org-Id", c.cloudOrgID)
	}

	return req, nil
}
