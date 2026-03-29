package tracker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ListOptions specifies the optional parameters for list methods that
// support page-based pagination.
type ListOptions struct {
	// Page is the page number to retrieve (1-based).
	Page int `url:"page,omitempty"`

	// PerPage is the number of items per page.
	PerPage int `url:"perPage,omitempty"`
}

// IssueSearchOptions specifies the optional query parameters for
// issue search requests.
type IssueSearchOptions struct {
	ListOptions

	// Expand specifies additional fields to include in the response
	// (e.g., "transitions").
	Expand string `url:"expand,omitempty"`
}

// ScrollSearchOptions specifies the optional query parameters for
// scroll-based issue search requests.
type ScrollSearchOptions struct {
	// ScrollType specifies the scroll type ("sorted" or "unsorted").
	ScrollType string `url:"scrollType,omitempty"`

	// PerScroll is the number of items per scroll page.
	PerScroll int `url:"perScroll,omitempty"`

	// ScrollTTLMs is the scroll cursor TTL in milliseconds.
	ScrollTTLMs int `url:"scrollTTLMillis,omitempty"`
}

// ChangelogOptions specifies the optional query parameters for
// retrieving issue changelog entries.
type ChangelogOptions struct {
	// ID is the changelog cursor ID for pagination.
	ID string `url:"id,omitempty"`

	// PerPage is the number of changelog entries per page.
	PerPage int `url:"perPage,omitempty"`
}

// IssueSearchRequest represents the request body for searching issues.
type IssueSearchRequest struct {
	// Filter is a map of field filters for the search.
	Filter map[string]any `json:"filter,omitempty"`

	// Query is a query string in the Yandex Tracker query language.
	Query *string `json:"query,omitempty"`

	// Order is the sort order for results.
	Order *string `json:"order,omitempty"`

	// Keys is a list of specific issue keys to search for.
	Keys []string `json:"keys,omitempty"`

	// Queue filters issues by queue key.
	Queue *string `json:"queue,omitempty"`
}

// IssueRequest represents the request body for creating or editing an issue.
type IssueRequest struct {
	// Summary is the issue title.
	Summary *string `json:"summary,omitempty"`

	// Description is the issue description.
	Description *string `json:"description,omitempty"`

	// Queue is the queue key to create the issue in.
	Queue *string `json:"queue,omitempty"`

	// Type is the issue type key.
	Type *string `json:"type,omitempty"`

	// Priority is the issue priority key.
	Priority *string `json:"priority,omitempty"`

	// Parent is the parent issue key.
	Parent *string `json:"parent,omitempty"`

	// Assignee is the assignee user ID.
	Assignee *string `json:"assignee,omitempty"`

	// Followers is a list of follower user IDs.
	Followers []string `json:"followers,omitempty"`

	// Tags is a list of issue tags.
	Tags []string `json:"tags,omitempty"`

	// Sprint is the sprint to assign the issue to.
	// Can be a sprint ID (string) or an object with ID.
	Sprint any `json:"sprint,omitempty"`

	// Deadline is the issue deadline in YYYY-MM-DD format.
	Deadline *string `json:"deadline,omitempty"`

	// Start is the planned start date in YYYY-MM-DD format.
	Start *string `json:"start,omitempty"`

	// End is the planned end date in YYYY-MM-DD format.
	End *string `json:"end,omitempty"`

	// Components is a list of component IDs or keys.
	Components []any `json:"components,omitempty"`

	// AffectedVersions is a list of affected version IDs.
	AffectedVersions []any `json:"affectedVersions,omitempty"`

	// FixVersions is a list of fix version IDs.
	FixVersions []any `json:"fixVersions,omitempty"`

	// StoryPoints is the story point estimate.
	StoryPoints *float64 `json:"storyPoints,omitempty"`

	// OriginalEstimation is the original time estimate in ISO 8601 duration format (e.g., "PT1H30M").
	OriginalEstimation *string `json:"originalEstimation,omitempty"`

	// Estimation is the current time estimate in ISO 8601 duration format (e.g., "PT1H30M").
	Estimation *string `json:"estimation,omitempty"`

	// Spent is the time spent in ISO 8601 duration format (e.g., "PT1H30M").
	Spent *string `json:"spent,omitempty"`
}

// TransitionRequest represents the request body for executing a
// status transition on an issue.
type TransitionRequest struct {
	// Comment is an optional comment to add when executing the transition.
	Comment *string `json:"comment,omitempty"`
}

// IssueMoveOptions specifies the query parameters for moving an issue
// to a different queue.
type IssueMoveOptions struct {
	// Queue is the target queue key (required).
	Queue string `url:"queue"`

	// MoveAllFields controls whether all fields are moved to the new queue.
	MoveAllFields *bool `url:"moveAllFields,omitempty"`

	// InitialStatus controls whether the issue is set to the initial status
	// of the target queue.
	InitialStatus *bool `url:"initialStatus,omitempty"`

	// ExpandTransitions controls whether available transitions are included
	// in the response.
	ExpandTransitions *bool `url:"expandTransitions,omitempty"`
}

// CommentListOptions specifies the optional parameters for listing comments.
// Comments use cursor-based pagination (not page-based).
type CommentListOptions struct {
	// ID is the comment ID cursor for pagination.
	// Results start after this comment ID.
	ID string `url:"id,omitempty"`

	// PerPage is the number of comments per page (default 50).
	PerPage int `url:"perPage,omitempty"`
}

// CommentRequest represents the request body for creating or editing a comment.
type CommentRequest struct {
	// Text is the comment text (required for create, required for edit).
	Text *string `json:"text,omitempty"`

	// AttachmentIDs is a list of attachment IDs to attach to the comment.
	AttachmentIDs []string `json:"attachmentIds,omitempty"`

	// Summonees is a list of user IDs or usernames to summon.
	Summonees []string `json:"summonees,omitempty"`

	// MaillistSummonees is a list of mailing list addresses to notify.
	MaillistSummonees []string `json:"maillistSummonees,omitempty"`
}

// LinkRequest represents the request body for creating a link between issues.
type LinkRequest struct {
	// Relationship is the link type (e.g., "relates", "is dependent by", "depends on").
	Relationship *string `json:"relationship,omitempty"`

	// Issue is the key of the issue to link to (e.g., "TREK-2").
	Issue *string `json:"issue,omitempty"`
}

// ChecklistItemRequest represents the request body for creating or editing a checklist item.
type ChecklistItemRequest struct {
	// Text is the checklist item text.
	Text *string `json:"text,omitempty"`

	// Checked is whether the item is checked.
	Checked *bool `json:"checked,omitempty"`

	// Assignee is the user ID to assign to the checklist item.
	Assignee *string `json:"assignee,omitempty"`

	// Deadline is the deadline for the checklist item.
	Deadline *ChecklistDeadline `json:"deadline,omitempty"`
}

// WorklogRequest represents the request body for creating or editing a worklog entry.
type WorklogRequest struct {
	// Start is the start date/time of the work.
	Start *Timestamp `json:"start,omitempty"`

	// Duration is the time spent in ISO 8601 format.
	Duration *Duration `json:"duration,omitempty"`

	// Comment is an optional text description of the work.
	Comment *string `json:"comment,omitempty"`
}

// QueueCreateRequest represents the request body for creating a queue.
type QueueCreateRequest struct {
	Key              *string            `json:"key,omitempty"`
	Name             *string            `json:"name,omitempty"`
	Lead             *string            `json:"lead,omitempty"`
	DefaultType      *string            `json:"defaultType,omitempty"`
	DefaultPriority  *string            `json:"defaultPriority,omitempty"`
	IssueTypesConfig []*IssueTypeConfig `json:"issueTypesConfig,omitempty"`
}

// IssueTypeConfig represents an issue type configuration within a queue.
type IssueTypeConfig struct {
	IssueType   *string  `json:"issueType,omitempty"`
	Workflow    *string  `json:"workflow,omitempty"`
	Resolutions []string `json:"resolutions,omitempty"`
}

// QueueGetOptions specifies optional parameters for getting a queue.
type QueueGetOptions struct {
	Expand string `url:"expand,omitempty"`
}

// QueueListOptions specifies optional parameters for listing queues.
type QueueListOptions struct {
	ListOptions

	Expand string `url:"expand,omitempty"`
}

// QueuePermissionsUpdateRequest represents the request body for updating queue permissions.
type QueuePermissionsUpdateRequest struct {
	Create *PermissionUpdate `json:"create,omitempty"`
	Write  *PermissionUpdate `json:"write,omitempty"`
	Read   *PermissionUpdate `json:"read,omitempty"`
	Grant  *PermissionUpdate `json:"grant,omitempty"`
}

// PermissionUpdate represents an add/remove operation for a permission group.
type PermissionUpdate struct {
	Add    *PermissionSubjects `json:"add,omitempty"`
	Remove *PermissionSubjects `json:"remove,omitempty"`
}

// PermissionSubjects represents the users, groups, and roles to add or remove.
type PermissionSubjects struct {
	Users  []string `json:"users,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

// FieldCreateRequest represents the request body for creating a field.
type FieldCreateRequest struct {
	ID              *string          `json:"id,omitempty"`
	Name            *FieldName       `json:"name,omitempty"`
	Category        *string          `json:"category,omitempty"`
	Type            *string          `json:"type,omitempty"`
	OptionsProvider *OptionsProvider `json:"optionsProvider,omitempty"`
	Order           *int             `json:"order,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Readonly        *bool            `json:"readonly,omitempty"`
	Visible         *bool            `json:"visible,omitempty"`
	Hidden          *bool            `json:"hidden,omitempty"`
	Container       *bool            `json:"container,omitempty"`
}

// FieldEditRequest represents the request body for editing a field.
type FieldEditRequest struct {
	Name            *FieldName       `json:"name,omitempty"`
	Category        *string          `json:"category,omitempty"`
	Order           *int             `json:"order,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Readonly        *bool            `json:"readonly,omitempty"`
	Visible         *bool            `json:"visible,omitempty"`
	Hidden          *bool            `json:"hidden,omitempty"`
	OptionsProvider *OptionsProvider `json:"optionsProvider,omitempty"`
}

// FieldEditOptions specifies optional parameters for editing a global field.
// Version is required for global field edits (optimistic locking).
type FieldEditOptions struct {
	Version int `url:"version"`
}

// ComponentRequest represents the request body for creating or editing a component.
type ComponentRequest struct {
	Name        *string `json:"name,omitempty"`
	Queue       *string `json:"queue,omitempty"`
	Description *string `json:"description,omitempty"`
	Lead        *string `json:"lead,omitempty"`
	AssignAuto  *bool   `json:"assignAuto,omitempty"`
}

// PriorityListOptions specifies the optional parameters for listing priorities.
type PriorityListOptions struct {
	Localized *bool `url:"localized,omitempty"`
}

// UserListOptions specifies the optional parameters for listing users.
type UserListOptions struct {
	ListOptions
}

// BulkMoveRequest represents the request body for bulk moving issues.
type BulkMoveRequest struct {
	Queue         *string        `json:"queue,omitempty"`
	Issues        []string       `json:"issues,omitempty"`
	Values        map[string]any `json:"values,omitempty"`
	MoveAllFields *bool          `json:"moveAllFields,omitempty"`
	InitialStatus *bool          `json:"initialStatus,omitempty"`
	Notify        *bool          `json:"notify,omitempty"`
}

// BulkUpdateRequest represents the request body for bulk updating issues.
type BulkUpdateRequest struct {
	Issues []string       `json:"issues,omitempty"`
	Values map[string]any `json:"values,omitempty"`
	Notify *bool          `json:"notify,omitempty"`
}

// BulkTransitionRequest represents the request body for bulk transitioning issues.
type BulkTransitionRequest struct {
	Transition *string        `json:"transition,omitempty"`
	Issues     []string       `json:"issues,omitempty"`
	Values     map[string]any `json:"values,omitempty"`
	Notify     *bool          `json:"notify,omitempty"`
}

// BoardCreateRequest represents the request body for creating a board.
type BoardCreateRequest struct {
	Name         *string        `json:"name,omitempty"`
	BoardType    *string        `json:"boardType,omitempty"`
	DefaultQueue *DefaultQueue  `json:"defaultQueue,omitempty"`
	Filter       map[string]any `json:"filter,omitempty"`
	OrderBy      *string        `json:"orderBy,omitempty"`
	OrderAsc     *bool          `json:"orderAsc,omitempty"`
	Query        *string        `json:"query,omitempty"`
	UseRanking   *bool          `json:"useRanking,omitempty"`
	Country      *Country       `json:"country,omitempty"`
}

// DefaultQueue identifies the default queue for a board by ID or key.
type DefaultQueue struct {
	ID  *string `json:"id,omitempty"`
	Key *string `json:"key,omitempty"`
}

// BoardEditRequest represents the request body for editing a board.
type BoardEditRequest struct {
	Name       *string             `json:"name,omitempty"`
	Columns    []*ColumnDefinition `json:"columns,omitempty"`
	Filter     map[string]any      `json:"filter,omitempty"`
	OrderBy    *string             `json:"orderBy,omitempty"`
	OrderAsc   *bool               `json:"orderAsc,omitempty"`
	Query      *string             `json:"query,omitempty"`
	UseRanking *bool               `json:"useRanking,omitempty"`
	Country    *Country            `json:"country,omitempty"`
}

// ColumnDefinition represents a column definition within a board edit request.
type ColumnDefinition struct {
	ID       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Statuses any     `json:"statuses,omitempty"`
}

// ColumnCreateRequest represents the request body for creating a column on a board.
type ColumnCreateRequest struct {
	Name     *string  `json:"name,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
}

// ColumnEditRequest represents the request body for editing a column on a board.
type ColumnEditRequest struct {
	Name     *string  `json:"name,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
}

// SprintCreateRequest represents the request body for creating a sprint.
// Note: Board ID is passed in the request body (not the URL).
// The API endpoint is POST /v3/sprints.
type SprintCreateRequest struct {
	Name      *string   `json:"name,omitempty"`
	Board     *BoardRef `json:"board,omitempty"`
	StartDate *string   `json:"startDate,omitempty"`
	EndDate   *string   `json:"endDate,omitempty"`
}

// EntityType represents the type of entity (project, portfolio, or goal).
type EntityType string

const (
	EntityTypeProject   EntityType = "project"
	EntityTypePortfolio EntityType = "portfolio"
	EntityTypeGoal      EntityType = "goal"
)

// EntityCreateRequest represents the request body for creating an entity.
type EntityCreateRequest struct {
	Fields *EntityCreateFields `json:"fields,omitempty"`
	Links  []*EntityLink       `json:"links,omitempty"`
}

// EntityCreateFields holds the fields for creating an entity.
type EntityCreateFields struct {
	Summary      *string       `json:"summary,omitempty"`
	TeamAccess   *bool         `json:"teamAccess,omitempty"`
	Description  *string       `json:"description,omitempty"`
	MarkupType   *string       `json:"markupType,omitempty"`
	Author       *string       `json:"author,omitempty"`
	Lead         *string       `json:"lead,omitempty"`
	TeamUsers    []string      `json:"teamUsers,omitempty"`
	Clients      []string      `json:"clients,omitempty"`
	Followers    []string      `json:"followers,omitempty"`
	Start        *string       `json:"start,omitempty"`
	End          *string       `json:"end,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	ParentEntity *ParentEntity `json:"parentEntity,omitempty"`
	EntityStatus *string       `json:"entityStatus,omitempty"`
}

// EntityLink represents a link in entity create/update requests.
type EntityLink struct {
	Relationship *string `json:"relationship,omitempty"`
	Entity       *string `json:"entity,omitempty"`
}

// EntityUpdateRequest represents the request body for updating an entity.
type EntityUpdateRequest struct {
	Fields  *EntityCreateFields `json:"fields,omitempty"`
	Comment *string             `json:"comment,omitempty"`
	Links   []*EntityLink       `json:"links,omitempty"`
}

// EntityGetOptions specifies optional query parameters for getting an entity.
type EntityGetOptions struct {
	Fields string `url:"fields,omitempty"`
	Expand string `url:"expand,omitempty"`
}

// EntityDeleteOptions specifies optional query parameters for deleting an entity.
type EntityDeleteOptions struct {
	WithBoard *bool `url:"withBoard,omitempty"`
}

// EntitySearchOptions specifies query parameters for entity search.
type EntitySearchOptions struct {
	ListOptions

	Fields string `url:"fields,omitempty"`
}

// EntitySearchRequest represents the POST body for searching entities.
type EntitySearchRequest struct {
	Input    *string        `json:"input,omitempty"`
	Filter   map[string]any `json:"filter,omitempty"`
	OrderBy  *string        `json:"orderBy,omitempty"`
	OrderAsc *bool          `json:"orderAsc,omitempty"`
	RootOnly *bool          `json:"rootOnly,omitempty"`
}

// EntityEventsOptions specifies query parameters for entity event history.
type EntityEventsOptions struct {
	PerPage        int    `url:"perPage,omitempty"`
	From           string `url:"from,omitempty"`
	Selected       string `url:"selected,omitempty"`
	NewEventsOnTop *bool  `url:"newEventsOnTop,omitempty"`
	Direction      string `url:"direction,omitempty"`
}

// EntityBulkChangeRequest represents the request body for bulk entity changes.
type EntityBulkChangeRequest struct {
	MetaEntities []string                `json:"metaEntities,omitempty"`
	Values       *EntityBulkChangeValues `json:"values,omitempty"`
}

// EntityBulkChangeValues represents the values to apply in a bulk change.
type EntityBulkChangeValues struct {
	Fields  map[string]any `json:"fields,omitempty"`
	Comment *string        `json:"comment,omitempty"`
	Links   []*EntityLink  `json:"links,omitempty"`
}

// EntityAccessUpdateRequest represents the request body for updating entity access.
type EntityAccessUpdateRequest struct {
	PermissionSources []string         `json:"permissionSources,omitempty"`
	ACL               *EntityACLUpdate `json:"acl,omitempty"`
}

// EntityACLUpdate represents the grant/revoke structure for access updates.
type EntityACLUpdate struct {
	Grant  *EntityACLPermissions `json:"grant,omitempty"`
	Revoke *EntityACLPermissions `json:"revoke,omitempty"`
}

// EntityACLPermissions represents permissions grouped by access level.
type EntityACLPermissions struct {
	Read  *EntityACLSubjects `json:"READ,omitempty"`
	Write *EntityACLSubjects `json:"WRITE,omitempty"`
	Grant *EntityACLSubjects `json:"GRANT,omitempty"`
}

// EntityACLSubjects represents the users, groups, and roles for a permission.
type EntityACLSubjects struct {
	Users  []string `json:"users,omitempty"`
	Groups []int    `json:"groups,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

// ChecklistMoveRequest represents the request body for moving a checklist item.
type ChecklistMoveRequest struct {
	// Before is the ID of the checklist item before which the moved item
	// will be inserted.
	Before *string `json:"before,omitempty"`
}

// EntityLinkDeleteOptions specifies the query parameters for deleting an entity link.
type EntityLinkDeleteOptions struct {
	// Right is the ID of the linked entity to unlink.
	Right string `url:"right"`
}

// EntityCommentCreateOptions specifies optional query parameters for
// creating or editing an entity comment.
type EntityCommentCreateOptions struct {
	// IsAddToFollowers controls whether the comment author is added as a follower.
	IsAddToFollowers *bool `url:"isAddToFollowers,omitempty"`

	// Notify controls whether users in Author, Lead, Participants, Customers,
	// and Followers fields are notified.
	Notify *bool `url:"notify,omitempty"`

	// NotifyAuthor controls whether the change author is notified.
	NotifyAuthor *bool `url:"notifyAuthor,omitempty"`
}

// TriggerCreateRequest represents the request body for creating a trigger.
type TriggerCreateRequest struct {
	Name       *string                `json:"name,omitempty"`
	Actions    []*AutomationAction    `json:"actions,omitempty"`
	Conditions []*AutomationCondition `json:"conditions,omitempty"`
	Active     *bool                  `json:"active,omitempty"`
}

// TriggerUpdateRequest represents the request body for updating a trigger.
type TriggerUpdateRequest struct {
	Name       *string                `json:"name,omitempty"`
	Actions    []*AutomationAction    `json:"actions,omitempty"`
	Conditions []*AutomationCondition `json:"conditions,omitempty"`
	Active     *bool                  `json:"active,omitempty"`
	Before     *int                   `json:"before,omitempty"`
}

// TriggerUpdateOptions specifies the query parameters for updating a trigger.
// Version is required for optimistic locking.
type TriggerUpdateOptions struct {
	Version int `url:"version"`
}

// AutoActionCreateRequest represents the request body for creating an auto-action.
type AutoActionCreateRequest struct {
	Name                *string             `json:"name,omitempty"`
	Filter              any                 `json:"filter,omitempty"`
	Query               *string             `json:"query,omitempty"`
	Actions             []*AutomationAction `json:"actions,omitempty"`
	Active              *bool               `json:"active,omitempty"`
	EnableNotifications *bool               `json:"enableNotifications,omitempty"`
	IntervalMillis      *int64              `json:"intervalMillis,omitempty"`
	Calendar            *AutoActionCalendar `json:"calendar,omitempty"`
}

// AutoActionUpdateRequest represents the request body for updating an auto-action.
type AutoActionUpdateRequest struct {
	Name                *string             `json:"name,omitempty"`
	Filter              any                 `json:"filter,omitempty"`
	Query               *string             `json:"query,omitempty"`
	Actions             []*AutomationAction `json:"actions,omitempty"`
	Active              *bool               `json:"active,omitempty"`
	EnableNotifications *bool               `json:"enableNotifications,omitempty"`
	IntervalMillis      *int64              `json:"intervalMillis,omitempty"`
	Calendar            *AutoActionCalendar `json:"calendar,omitempty"`
}

// MacroCreateRequest represents the request body for creating a macro in a queue.
type MacroCreateRequest struct {
	Name        *string `json:"name,omitempty"`
	Body        *string `json:"body,omitempty"`
	IssueUpdate any     `json:"issueUpdate,omitempty"`
}

// MacroEditRequest represents the request body for editing a macro in a queue.
type MacroEditRequest struct {
	Name        *string `json:"name,omitempty"`
	Body        *string `json:"body,omitempty"`
	IssueUpdate any     `json:"issueUpdate,omitempty"`
}

// ExternalLinkCreateRequest represents the body for creating an external link on an issue.
type ExternalLinkCreateRequest struct {
	Relationship *string `json:"relationship,omitempty"`
	Key          *string `json:"key,omitempty"`
	Origin       *string `json:"origin,omitempty"`
}

// ExternalLinkCreateOptions specifies the optional query parameters for creating an external link.
type ExternalLinkCreateOptions struct {
	Backlink *bool `url:"backlink,omitempty"`
}

// ImportIssueRequest represents the body for importing an issue with preserved timestamps.
// The queue and summary fields are required; createdAt and createdBy preserve the original creation metadata.
//
// Timestamps (createdAt, updatedAt, resolvedAt) are *string in "YYYY-MM-DDThh:mm:ss.sss+hhmm" format
// because the consumer controls the exact format string for import operations.
//
// Constraints enforced by the API (not this library):
//   - createdAt must not exceed current time
//   - updatedAt must fall between createdAt and current time
//   - resolvedAt must fall between createdAt and updatedAt
//   - status must exist in the queue's workflow for the selected issue type
//   - user must have issue creation permissions in target queue
type ImportIssueRequest struct {
	Queue              *string  `json:"queue,omitempty"`
	Summary            *string  `json:"summary,omitempty"`
	Key                *string  `json:"key,omitempty"`
	CreatedAt          *string  `json:"createdAt,omitempty"`
	CreatedBy          *string  `json:"createdBy,omitempty"`
	UpdatedAt          *string  `json:"updatedAt,omitempty"`
	UpdatedBy          *string  `json:"updatedBy,omitempty"`
	ResolvedAt         *string  `json:"resolvedAt,omitempty"`
	ResolvedBy         *string  `json:"resolvedBy,omitempty"`
	Status             *int     `json:"status,omitempty"`
	Deadline           *string  `json:"deadline,omitempty"`
	Resolution         *int     `json:"resolution,omitempty"`
	Type               *int     `json:"type,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Start              *string  `json:"start,omitempty"`
	End                *string  `json:"end,omitempty"`
	Assignee           *string  `json:"assignee,omitempty"`
	Priority           *int     `json:"priority,omitempty"`
	AffectedVersions   []int    `json:"affectedVersions,omitempty"`
	FixVersions        []int    `json:"fixVersions,omitempty"`
	Components         []int    `json:"components,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	Sprint             []int    `json:"sprint,omitempty"`
	Followers          []string `json:"followers,omitempty"`
	Access             []string `json:"access,omitempty"`
	Unique             *string  `json:"unique,omitempty"`
	FollowingMaillists []string `json:"followingMaillists,omitempty"`
	OriginalEstimation *int     `json:"originalEstimation,omitempty"`
	Estimation         *int     `json:"estimation,omitempty"`
	Spent              *int     `json:"spent,omitempty"`
	StoryPoints        *float64 `json:"storyPoints,omitempty"`
	VotedBy            []string `json:"votedBy,omitempty"`
	FavoritedBy        []string `json:"favoritedBy,omitempty"`
}

// ImportCommentRequest represents the body for importing a comment with preserved timestamps.
//
// Constraints enforced by the API (not this library):
//   - comments must be imported in chronological order
//   - createdAt must not exceed current time
//   - updatedAt must fall between createdAt and current time
type ImportCommentRequest struct {
	Text      *string `json:"text,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	CreatedBy *string `json:"createdBy,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	UpdatedBy *string `json:"updatedBy,omitempty"`
}

// ImportLinkRequest represents the body for importing a link with preserved timestamps.
type ImportLinkRequest struct {
	Relationship *string `json:"relationship,omitempty"`
	Issue        *string `json:"issue,omitempty"`
	CreatedAt    *string `json:"createdAt,omitempty"`
	CreatedBy    *string `json:"createdBy,omitempty"`
	UpdatedAt    *string `json:"updatedAt,omitempty"`
	UpdatedBy    *string `json:"updatedBy,omitempty"`
}

// ImportFileOptions specifies the required query parameters for file import.
// These are sent as URL query parameters, not in the request body.
// The request body is multipart file data.
type ImportFileOptions struct {
	Filename  string `url:"filename"`
	CreatedAt string `url:"createdAt"`
	CreatedBy string `url:"createdBy"`
}

// DashboardCreateRequest represents the body for creating a dashboard.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/dashboards/create-dashboard
type DashboardCreateRequest struct {
	Name   *string         `json:"name,omitempty"`
	Layout *string         `json:"layout,omitempty"`
	Owner  *DashboardOwner `json:"owner,omitempty"`
}

// DashboardOwner identifies a dashboard owner in create/update requests.
// The request takes {"id": "username"}, which differs from the response
// that returns a full User object.
type DashboardOwner struct {
	ID *string `json:"id,omitempty"`
}

// WidgetCreateRequest represents the body for creating a widget on a dashboard.
// The widget type is specified as a separate parameter in CreateWidget (used to
// construct the URL path), not in the request body.
//
// Description is the only required field. All remaining widget-type-specific
// parameters should be set in the Parameters map.
//
// Yandex Tracker API docs: https://yandex.ru/support/tracker/en/api-ref/dashboards/create-widget
type WidgetCreateRequest struct {
	Description *string        `json:"-"`
	Parameters  map[string]any `json:"-"`
}

// widgetCreateRequestKnownKeys lists the JSON keys that map to known
// WidgetCreateRequest struct fields.
var widgetCreateRequestKnownKeys = []string{"description"}

// MarshalJSON implements the json.Marshaler interface for WidgetCreateRequest.
// It builds a flat JSON object containing the description field and
// all entries from the Parameters map.
func (r WidgetCreateRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	for k, v := range r.Parameters {
		m[k] = v
	}
	if r.Description != nil {
		m["description"] = *r.Description
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements the json.Unmarshaler interface for WidgetCreateRequest.
// It decodes the description field into the struct field and stores
// any remaining keys in the Parameters map.
func (r *WidgetCreateRequest) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["description"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			r.Description = &s
		}
	}

	for _, k := range widgetCreateRequestKnownKeys {
		delete(raw, k)
	}
	if len(raw) > 0 {
		r.Parameters = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of parameter
			r.Parameters[k] = val
		}
	}

	return nil
}

// addOptions adds the parameters in opts as URL query parameters to s.
// opts must be a struct or a pointer to a struct whose fields have
// "url" struct tags.
func addOptions(s string, opts any) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return s, nil
		}
		v = v.Elem()
	}

	qs := u.Query()
	encodeFields(v, qs)

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// encodeFields iterates over struct fields and encodes them as query parameters.
// It handles embedded structs by recursing into them.
func encodeFields(v reflect.Value, qs url.Values) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		// Handle embedded structs by recursing.
		if structField.Anonymous && field.Kind() == reflect.Struct {
			encodeFields(field, qs)
			continue
		}

		tag := structField.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}

		name, tagOpts := parseTag(tag)

		// Handle pointer fields: a nil pointer is always omitted;
		// a non-nil pointer means the user intentionally set the value,
		// so we skip the omitempty check on the dereferenced value.
		isPtr := field.Kind() == reflect.Ptr
		if isPtr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		if !isPtr && tagOpts == "omitempty" && isEmptyValue(field) {
			continue
		}

		qs.Set(name, fmt.Sprint(field.Interface()))
	}
}

// parseTag splits a struct tag into its name and options.
// A tag "name,opt" returns ("name", "opt").
func parseTag(tag string) (string, string) {
	name, opts, _ := strings.Cut(tag, ",")
	return name, opts
}

// isEmptyValue checks if a reflect.Value is the zero value for its type.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}
