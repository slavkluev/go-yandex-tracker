package tracker

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// setup creates a test Client and ServeMux connected to an httptest.Server.
// The server is automatically closed when the test finishes.
func setup(t *testing.T) (client *Client, mux *http.ServeMux) {
	t.Helper()
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	client = NewClient(
		WithBaseURL(server.URL+"/"),
		WithOAuthToken("test-token"),
		WithOrgID("test-org"),
	)
	return client, mux
}

// testMethod verifies the HTTP method of a request.
func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

// testHeader verifies a specific header value in a request.
func testHeader(t *testing.T, r *http.Request, header, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) = %q, want %q", header, got, want)
	}
}

// testBody verifies the body content of a request.
func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}
	if got := strings.TrimSpace(string(b)); got != want {
		t.Errorf("Request body = %s, want %s", got, want)
	}
}

func TestNewClient_DefaultBaseURL(t *testing.T) {
	c := NewClient()
	want := "https://api.tracker.yandex.net/"
	if got := c.baseURL.String(); got != want {
		t.Errorf("NewClient().baseURL = %q, want %q", got, want)
	}
}

func TestNewClient_DefaultHTTPClient(t *testing.T) {
	c := NewClient()
	if c.client == nil {
		t.Error("NewClient().client is nil, want non-nil *http.Client")
	}
}

func TestNewClient_WithOAuthToken(t *testing.T) {
	c := NewClient(WithOAuthToken("my-oauth-token"))
	if c.token != "my-oauth-token" {
		t.Errorf("token = %q, want %q", c.token, "my-oauth-token")
	}
	if c.authType != authOAuth {
		t.Errorf("authType = %v, want %v", c.authType, authOAuth)
	}
}

func TestNewClient_WithIAMToken(t *testing.T) {
	c := NewClient(WithIAMToken("my-iam-token"))
	if c.token != "my-iam-token" {
		t.Errorf("token = %q, want %q", c.token, "my-iam-token")
	}
	if c.authType != authIAM {
		t.Errorf("authType = %v, want %v", c.authType, authIAM)
	}
}

func TestNewClient_WithOrgID(t *testing.T) {
	c := NewClient(WithOrgID("org123"))
	if c.orgID != "org123" {
		t.Errorf("orgID = %q, want %q", c.orgID, "org123")
	}
}

func TestNewClient_WithCloudOrgID(t *testing.T) {
	c := NewClient(WithCloudOrgID("cloud456"))
	if c.cloudOrgID != "cloud456" {
		t.Errorf("cloudOrgID = %q, want %q", c.cloudOrgID, "cloud456")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := NewClient(WithHTTPClient(custom))
	if c.client != custom {
		t.Error("NewClient(WithHTTPClient(custom)).client != custom")
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	c := NewClient(WithBaseURL("http://custom.example.com/api/"))
	want := "http://custom.example.com/api/"
	if got := c.baseURL.String(); got != want {
		t.Errorf("baseURL = %q, want %q", got, want)
	}
}

func TestNewClient_IssuesServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Issues == nil {
		t.Error("NewClient().Issues is nil, want non-nil *IssuesService")
	}
}

func TestNewClient_QueuesServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Queues == nil {
		t.Error("NewClient().Queues is nil, want non-nil *QueuesService")
	}
}

func TestNewClient_FieldsServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Fields == nil {
		t.Error("NewClient().Fields is nil, want non-nil *FieldsService")
	}
}

func TestNewClient_ComponentsServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Components == nil {
		t.Error("NewClient().Components is nil, want non-nil *ComponentsService")
	}
}

func TestNewClient_IssueTypesServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.IssueTypes == nil {
		t.Error("NewClient().IssueTypes is nil, want non-nil *IssueTypesService")
	}
}

func TestNewClient_StatusesServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Statuses == nil {
		t.Error("NewClient().Statuses is nil, want non-nil *StatusesService")
	}
}

func TestNewClient_ResolutionsServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Resolutions == nil {
		t.Error("NewClient().Resolutions is nil, want non-nil *ResolutionsService")
	}
}

func TestNewClient_PrioritiesServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Priorities == nil {
		t.Error("NewClient().Priorities is nil, want non-nil *PrioritiesService")
	}
}

func TestNewClient_UsersServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.Users == nil {
		t.Error("NewClient().Users is nil, want non-nil *UsersService")
	}
}

func TestNewClient_BulkChangeServiceNotNil(t *testing.T) {
	c := NewClient()
	if c.BulkChange == nil {
		t.Error("NewClient().BulkChange is nil, want non-nil *BulkChangeService")
	}
}

func TestNewRequest_GET_NilBody(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	// Check URL
	wantURL := defaultBaseURL + "v3/issues/QUEUE-1"
	if got := req.URL.String(); got != wantURL {
		t.Errorf("URL = %q, want %q", got, wantURL)
	}

	// No body
	if req.Body != nil {
		t.Error("Request body is non-nil for nil body input")
	}

	// No Content-Type header
	if got := req.Header.Get("Content-Type"); got != "" {
		t.Errorf("Content-Type = %q, want empty for nil body", got)
	}
}

func TestNewRequest_POST_WithBody(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	body := &IssueRequest{
		Summary: Ptr("Test issue"),
		Queue:   Ptr("TEST"),
	}
	req, err := c.NewRequest("POST", "v3/issues/", body)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	// Content-Type should be set
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	// Body should contain JSON
	b, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("Body is not valid JSON: %v", err)
	}
	if decoded["summary"] != "Test issue" {
		t.Errorf("body summary = %v, want %q", decoded["summary"], "Test issue")
	}
	if decoded["queue"] != "TEST" {
		t.Errorf("body queue = %v, want %q", decoded["queue"], "TEST")
	}
}

func TestNewRequest_OAuthHeader(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	want := "OAuth test-token"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestNewRequest_IAMHeader(t *testing.T) {
	c := NewClient(WithIAMToken("test-iam-token"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	want := "Bearer test-iam-token"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestNewRequest_OrgIDHeader(t *testing.T) {
	c := NewClient(WithOrgID("org123"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	want := "org123"
	if got := req.Header.Get("X-Org-Id"); got != want {
		t.Errorf("X-Org-Id = %q, want %q", got, want)
	}
}

func TestNewRequest_CloudOrgIDHeader(t *testing.T) {
	c := NewClient(WithCloudOrgID("cloud456"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	want := "cloud456"
	if got := req.Header.Get("X-Cloud-Org-Id"); got != want {
		t.Errorf("X-Cloud-Org-Id = %q, want %q", got, want)
	}
}

func TestNewRequest_BothOrgHeaders(t *testing.T) {
	c := NewClient(WithOrgID("org123"), WithCloudOrgID("cloud456"))
	req, err := c.NewRequest("GET", "v3/issues/QUEUE-1", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	if got := req.Header.Get("X-Org-Id"); got != "org123" {
		t.Errorf("X-Org-Id = %q, want %q", got, "org123")
	}
	if got := req.Header.Get("X-Cloud-Org-Id"); got != "cloud456" {
		t.Errorf("X-Cloud-Org-Id = %q, want %q", got, "cloud456")
	}
}

func TestDo_200_DecodesJSON(t *testing.T) {
	client, mux := setup(t)

	type testResource struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResource{Key: "TEST-1", Name: "Test"})
	})

	req, _ := client.NewRequest("GET", "v3/test", nil)
	var got testResource
	resp, err := client.Do(context.Background(), req, &got)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	want := testResource{Key: "TEST-1", Name: "Test"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Do decoded = %+v, want %+v", got, want)
	}
}

func TestDo_404_ReturnsErrorResponse(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errorMessages": []string{"Object not found"},
			"errors":        map[string]string{},
		})
	})

	req, _ := client.NewRequest("GET", "v3/test", nil)
	resp, err := client.Do(context.Background(), req, nil)
	if err == nil {
		t.Fatal("Do expected error for 404, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("error type = %T, want *ErrorResponse", err)
	}
	if errResp.Response.StatusCode != http.StatusNotFound {
		t.Errorf("ErrorResponse.StatusCode = %d, want %d", errResp.Response.StatusCode, http.StatusNotFound)
	}
}

func TestDo_NilV_200_NoError(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("DELETE /v3/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req, _ := client.NewRequest("DELETE", "v3/test", nil)
	resp, err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestDo_EmptyBody_NilV_NoError(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("POST /v3/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req, _ := client.NewRequest("POST", "v3/test", nil)
	resp, err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestDo_PassesContext(t *testing.T) {
	client, mux := setup(t)

	type contextKey string
	var gotCtx bool

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		// The context is set on the request; we can verify the request has it
		gotCtx = true
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	req, _ := client.NewRequest("GET", "v3/test", nil)
	_, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !gotCtx {
		t.Error("Handler was not called, context may not have been passed")
	}
}

func TestDo_CancelledContext(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req, _ := client.NewRequest("GET", "v3/test", nil)
	_, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Do expected error for cancelled context, got nil")
	}
}

func TestBareDo_200_ReturnsRawResponse(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("raw body content"))
	})

	req, _ := client.NewRequest("GET", "v3/test", nil)
	resp, err := client.BareDo(context.Background(), req)
	if err != nil {
		t.Fatalf("BareDo returned error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}
	if got := string(body); got != "raw body content" {
		t.Errorf("Body = %q, want %q", got, "raw body content")
	}
}

func TestBareDo_404_ReturnsError(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("GET /v3/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	req, _ := client.NewRequest("GET", "v3/test", nil)
	_, err := client.BareDo(context.Background(), req)
	if err == nil {
		t.Fatal("BareDo expected error for 404, got nil")
	}
}

func TestNewUploadRequest_CreatesMultipartRequest(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	fileContent := "hello file content"
	req, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader(fileContent))
	if err != nil {
		t.Fatalf("NewUploadRequest returned error: %v", err)
	}

	// Content-Type must start with multipart/form-data
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		t.Errorf("Content-Type = %q, want prefix %q", ct, "multipart/form-data")
	}

	// Parse multipart form to verify file field
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("ParseMultipartForm returned error: %v", err)
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("FormFile(\"file\") returned error: %v", err)
	}
	defer file.Close()

	if header.Filename != "test.txt" {
		t.Errorf("filename = %q, want %q", header.Filename, "test.txt")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Error reading file content: %v", err)
	}
	if got := string(content); got != fileContent {
		t.Errorf("file content = %q, want %q", got, fileContent)
	}
}

func TestNewUploadRequest_SetsAuthHeaders_OAuth(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	req, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("NewUploadRequest returned error: %v", err)
	}
	want := "OAuth test-token"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestNewUploadRequest_SetsAuthHeaders_IAM(t *testing.T) {
	c := NewClient(WithIAMToken("test-iam-token"))
	req, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("NewUploadRequest returned error: %v", err)
	}
	want := "Bearer test-iam-token"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestNewUploadRequest_SetsOrgHeaders(t *testing.T) {
	c := NewClient(
		WithOAuthToken("test-token"),
		WithOrgID("org123"),
		WithCloudOrgID("cloud456"),
	)
	req, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("NewUploadRequest returned error: %v", err)
	}
	if got := req.Header.Get("X-Org-Id"); got != "org123" {
		t.Errorf("X-Org-Id = %q, want %q", got, "org123")
	}
	if got := req.Header.Get("X-Cloud-Org-Id"); got != "cloud456" {
		t.Errorf("X-Cloud-Org-Id = %q, want %q", got, "cloud456")
	}
}

func TestNewUploadRequest_ResolvesURL(t *testing.T) {
	c := NewClient(WithOAuthToken("test-token"))
	req, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader("data"))
	if err != nil {
		t.Fatalf("NewUploadRequest returned error: %v", err)
	}
	wantURL := defaultBaseURL + "v3/issues/QUEUE-1/attachments"
	if got := req.URL.String(); got != wantURL {
		t.Errorf("URL = %q, want %q", got, wantURL)
	}
}

func TestNewUploadRequest_NoTrailingSlash(t *testing.T) {
	c := NewClient(
		WithBaseURL("http://example.com/api"),
		WithOAuthToken("test-token"),
	)
	_, err := c.NewUploadRequest("POST", "v3/issues/QUEUE-1/attachments", "test.txt", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error for baseURL without trailing slash, got nil")
	}
}
