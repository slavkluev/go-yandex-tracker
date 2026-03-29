package tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000-0700"

// Timestamp represents a time value in Yandex Tracker format.
// Yandex Tracker uses the format "2017-06-11T05:16:01.339+0000"
// which differs from RFC 3339 (uses +0000 instead of +00:00).
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses the Yandex Tracker timestamp format, falling back to RFC 3339.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}
	parsed, err := time.Parse(timestampLayout, str)
	if err != nil {
		// Fallback to RFC 3339
		parsed, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// It formats the time using the Yandex Tracker timestamp layout.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(timestampLayout) + `"`), nil
}

// Response wraps the standard http.Response and adds Yandex Tracker
// pagination and rate limiting fields parsed from response headers.
type Response struct {
	*http.Response

	// TotalCount is the total number of results, parsed from X-Total-Count header.
	TotalCount int

	// TotalPages is the total number of pages, parsed from X-Total-Pages header.
	TotalPages int

	// ScrollID is the scroll cursor ID for scroll-based pagination,
	// copied from the X-Scroll-Id header.
	ScrollID string

	// ScrollToken is the scroll token for scroll-based pagination,
	// copied from the X-Scroll-Token header.
	ScrollToken string

	// RetryAfter is the number of seconds to wait before retrying,
	// parsed from the Retry-After header on 429 responses.
	RetryAfter int
}

// newResponse creates a new Response from an http.Response, parsing
// pagination and rate limiting headers.
func newResponse(r *http.Response) *Response {
	resp := &Response{Response: r}

	if v := r.Header.Get("X-Total-Count"); v != "" {
		resp.TotalCount, _ = strconv.Atoi(v)
	}
	if v := r.Header.Get("X-Total-Pages"); v != "" {
		resp.TotalPages, _ = strconv.Atoi(v)
	}
	resp.ScrollID = r.Header.Get("X-Scroll-Id")
	resp.ScrollToken = r.Header.Get("X-Scroll-Token")
	if v := r.Header.Get("Retry-After"); v != "" {
		resp.RetryAfter, _ = strconv.Atoi(v)
	}

	return resp
}

// User represents a Yandex Tracker user.
// When embedded in other resources (e.g., Issue), only Self, ID, and Display
// are populated. When returned as a full resource (from GET /v3/myself,
// /v3/users), all fields are populated.
// Note: User.ID (*string, json:"id") is for embedded references.
// User.UID (*int, json:"uid") is for full responses. These are different JSON keys.
type User struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`

	// Full resource fields (from GET /v3/myself, /v3/users)
	UID                  *int       `json:"uid,omitempty"`
	Login                *string    `json:"login,omitempty"`
	TrackerUID           *int       `json:"trackerUid,omitempty"`
	PassportUID          *int       `json:"passportUid,omitempty"`
	CloudUID             *string    `json:"cloudUid,omitempty"`
	FirstName            *string    `json:"firstName,omitempty"`
	LastName             *string    `json:"lastName,omitempty"`
	Email                *string    `json:"email,omitempty"`
	External             *bool      `json:"external,omitempty"`
	HasLicense           *bool      `json:"hasLicense,omitempty"`
	Dismissed            *bool      `json:"dismissed,omitempty"`
	UseNewFilters        *bool      `json:"useNewFilters,omitempty"`
	DisableNotifications *bool      `json:"disableNotifications,omitempty"`
	FirstLoginDate       *Timestamp `json:"firstLoginDate,omitempty"`
	LastLoginDate        *Timestamp `json:"lastLoginDate,omitempty"`
	WelcomeMailSent      *bool      `json:"welcomeMailSent,omitempty"`
}

// Status represents an issue status in Yandex Tracker.
// When returned as a full resource (from GET /v3/statuses), all fields are
// populated. When embedded in other resources, only Self, ID, Key, and
// Display are populated.
type Status struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Key     *string `json:"key,omitempty"`
	Display *string `json:"display,omitempty"`

	// Full resource fields (from GET /v3/statuses)
	Version     *int    `json:"version,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Type        *string `json:"type,omitempty"`
}

// Queue represents a Yandex Tracker queue.
// When returned as a full resource, all fields are populated.
// When embedded in other resources (e.g., Issue), only Self, ID, Key,
// and Display are populated; other fields are nil.
type Queue struct {
	Self            *string    `json:"self,omitempty"`
	ID              *string    `json:"id,omitempty"`
	Key             *string    `json:"key,omitempty"`
	Display         *string    `json:"display,omitempty"`
	Version         *int       `json:"version,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Lead            *User      `json:"lead,omitempty"`
	AssignAuto      *bool      `json:"assignAuto,omitempty"`
	AllowExternals  *bool      `json:"allowExternals,omitempty"`
	DenyVoting      *bool      `json:"denyVoting,omitempty"`
	DefaultType     *IssueType `json:"defaultType,omitempty"`
	DefaultPriority *Priority  `json:"defaultPriority,omitempty"`
}

// Priority represents an issue priority in Yandex Tracker.
// When returned as a full resource (from GET /v3/priorities), all fields are
// populated. When embedded in other resources, only Self, ID, Key, and
// Display are populated.
type Priority struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Key     *string `json:"key,omitempty"`
	Display *string `json:"display,omitempty"`

	// Full resource fields (from GET /v3/priorities)
	Version *int    `json:"version,omitempty"`
	Name    *string `json:"name,omitempty"`
	Order   *int    `json:"order,omitempty"`
}

// IssueType represents the type of a Yandex Tracker issue.
// When returned as a full resource (from GET /v3/issuetypes), all fields are
// populated. When embedded in other resources, only Self, ID, Key, and
// Display are populated.
type IssueType struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Key     *string `json:"key,omitempty"`
	Display *string `json:"display,omitempty"`

	// Full resource fields (from GET /v3/issuetypes)
	Version     *int    `json:"version,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Deleted     *bool   `json:"deleted,omitempty"`
}

// Resolution represents the resolution of a Yandex Tracker issue.
// When returned as a full resource (from GET /v3/resolutions), all fields are
// populated. When embedded in other resources, only Self, ID, Key, and
// Display are populated.
type Resolution struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Key     *string `json:"key,omitempty"`
	Display *string `json:"display,omitempty"`

	// Full resource fields (from GET /v3/resolutions)
	Version     *int    `json:"version,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Order       *int    `json:"order,omitempty"`
}

// Issue represents a Yandex Tracker issue.
// Known API fields are mapped to typed struct fields.
// Any additional fields (custom queue fields) are captured in CustomFields.
type Issue struct {
	Self           *string          `json:"self,omitempty"`
	ID             *string          `json:"id,omitempty"`
	Key            *string          `json:"key,omitempty"`
	Summary        *string          `json:"summary,omitempty"`
	Description    *string          `json:"description,omitempty"`
	Type           *IssueType       `json:"type,omitempty"`
	Priority       *Priority        `json:"priority,omitempty"`
	Status         *Status          `json:"status,omitempty"`
	Queue          *Queue           `json:"queue,omitempty"`
	Assignee       *User            `json:"assignee,omitempty"`
	Followers      []*User          `json:"followers,omitempty"`
	CreatedBy      *User            `json:"createdBy,omitempty"`
	UpdatedBy      *User            `json:"updatedBy,omitempty"`
	CreatedAt      *Timestamp       `json:"createdAt,omitempty"`
	UpdatedAt      *Timestamp       `json:"updatedAt,omitempty"`
	ResolvedAt     *Timestamp       `json:"resolvedAt,omitempty"`
	Resolution     *Resolution      `json:"resolution,omitempty"`
	ChecklistItems []*ChecklistItem `json:"checklistItems,omitempty"`
	ChecklistDone  *int             `json:"checklistDone,omitempty"`
	ChecklistTotal *int             `json:"checklistTotal,omitempty"`

	// CustomFields holds any fields from the JSON response that do not
	// correspond to known Issue struct fields. These are typically
	// user-defined custom fields on a queue.
	CustomFields map[string]any `json:"-"`
}

// issueKnownKeys lists the JSON keys that map to known Issue struct fields.
// Keys not in this list are treated as custom fields.
var issueKnownKeys = []string{
	"self", "id", "key", "summary", "description",
	"type", "priority", "status", "queue",
	"assignee", "followers", "createdBy", "updatedBy",
	"createdAt", "updatedAt", "resolvedAt", "resolution",
	"checklistItems", "checklistDone", "checklistTotal",
}

// UnmarshalJSON implements the json.Unmarshaler interface for Issue.
// It decodes known fields into the struct and captures any remaining
// keys into the CustomFields map.
func (i *Issue) UnmarshalJSON(data []byte) error {
	// Use type alias to avoid infinite recursion.
	type IssueAlias Issue
	var alias IssueAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*i = Issue(alias)

	// Unmarshal the full JSON into a raw map to identify custom fields.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Remove known keys; remaining are custom fields.
	for _, knownKey := range issueKnownKeys {
		delete(raw, knownKey)
	}

	if len(raw) > 0 {
		i.CustomFields = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of custom field
			i.CustomFields[k] = val
		}
	}

	return nil
}

// Transition represents an issue status transition in Yandex Tracker.
type Transition struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
	To      *Status `json:"to,omitempty"`
}

// Changelog represents a change history entry for an issue.
type Changelog struct {
	Self      *string           `json:"self,omitempty"`
	ID        *string           `json:"id,omitempty"`
	Issue     *Issue            `json:"issue,omitempty"`
	UpdatedAt *Timestamp        `json:"updatedAt,omitempty"`
	UpdatedBy *User             `json:"updatedBy,omitempty"`
	Type      *string           `json:"type,omitempty"`
	Fields    []*ChangelogEvent `json:"fields,omitempty"`
}

// ChangelogEvent represents a single field change within a changelog entry.
type ChangelogEvent struct {
	Field *string `json:"field,omitempty"`
	From  any     `json:"from,omitempty"`
	To    any     `json:"to,omitempty"`
}

// IssueLink represents a link between two issues.
type IssueLink struct {
	Self      *string        `json:"self,omitempty"`
	ID        *string        `json:"id,omitempty"`
	Type      *IssueLinkType `json:"type,omitempty"`
	Direction *string        `json:"direction,omitempty"`
	Object    *Issue         `json:"object,omitempty"`
	CreatedBy *User          `json:"createdBy,omitempty"`
	UpdatedBy *User          `json:"updatedBy,omitempty"`
	CreatedAt *Timestamp     `json:"createdAt,omitempty"`
	UpdatedAt *Timestamp     `json:"updatedAt,omitempty"`
}

// IssueLinkType represents the type of a link between issues.
type IssueLinkType struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Inward  *string `json:"inward,omitempty"`
	Outward *string `json:"outward,omitempty"`
}

// Duration represents an ISO 8601 duration value (e.g., "PT1H30M", "P5D").
// It wraps time.Duration and handles JSON marshal/unmarshal in ISO 8601 format.
type Duration struct {
	time.Duration
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses ISO 8601 duration strings like "PT1H30M", "P5D", "P1W", "P1DT2H3M".
func (d *Duration) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		return fmt.Errorf("empty duration string")
	}
	if !strings.HasPrefix(str, "P") {
		return fmt.Errorf("invalid ISO 8601 duration: %q", str)
	}

	str = str[1:] // Strip leading "P"

	var total time.Duration

	// Split on "T" to separate date part and time part.
	datePart := str
	timePart := ""
	if idx := strings.Index(str, "T"); idx >= 0 {
		datePart = str[:idx]
		timePart = str[idx+1:]
	}

	// Parse date part: Y, M (months), W, D
	if datePart != "" {
		var err error
		total, err = parseDatePart(datePart)
		if err != nil {
			return fmt.Errorf("invalid ISO 8601 duration date part: %w", err)
		}
	}

	// Parse time part: H, M (minutes), S
	if timePart != "" {
		timeDur, err := parseTimePart(timePart)
		if err != nil {
			return fmt.Errorf("invalid ISO 8601 duration time part: %w", err)
		}
		total += timeDur
	}

	d.Duration = total
	return nil
}

// parseDatePart parses the date portion of an ISO 8601 duration (before "T").
// Handles Y (years, as 365 days), M (months, as 30 days), W (weeks), D (days).
func parseDatePart(s string) (time.Duration, error) {
	var total time.Duration
	for s != "" {
		// Find the next letter.
		i := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if i == 0 || i >= len(s) {
			return 0, fmt.Errorf("unexpected format: %q", s)
		}
		n, err := strconv.Atoi(s[:i])
		if err != nil {
			return 0, err
		}
		unit := s[i]
		s = s[i+1:]

		switch unit {
		case 'Y':
			total += time.Duration(n) * 365 * 24 * time.Hour
		case 'M':
			total += time.Duration(n) * 30 * 24 * time.Hour
		case 'W':
			total += time.Duration(n) * 7 * 24 * time.Hour
		case 'D':
			total += time.Duration(n) * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unknown date unit: %c", unit)
		}
	}
	return total, nil
}

// parseTimePart parses the time portion of an ISO 8601 duration (after "T").
// Handles H (hours), M (minutes), S (seconds).
func parseTimePart(s string) (time.Duration, error) {
	var total time.Duration
	for s != "" {
		i := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if i == 0 || i >= len(s) {
			return 0, fmt.Errorf("unexpected format: %q", s)
		}
		n, err := strconv.Atoi(s[:i])
		if err != nil {
			return 0, err
		}
		unit := s[i]
		s = s[i+1:]

		switch unit {
		case 'H':
			total += time.Duration(n) * time.Hour
		case 'M':
			total += time.Duration(n) * time.Minute
		case 'S':
			total += time.Duration(n) * time.Second
		default:
			return 0, fmt.Errorf("unknown time unit: %c", unit)
		}
	}
	return total, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It converts the duration to ISO 8601 format (e.g., "PT1H30M", "P5D", "P1W").
func (d Duration) MarshalJSON() ([]byte, error) {
	dur := d.Duration
	if dur == 0 {
		return []byte(`"PT0S"`), nil
	}

	var b strings.Builder
	b.WriteString("P")

	totalSeconds := int64(dur / time.Second)
	days := totalSeconds / (24 * 3600)
	remainder := totalSeconds % (24 * 3600)
	hours := remainder / 3600
	remainder %= 3600
	minutes := remainder / 60
	seconds := remainder % 60

	hasDatePart := false
	if days > 0 {
		if days%7 == 0 {
			b.WriteString(strconv.FormatInt(days/7, 10))
			b.WriteByte('W')
		} else {
			b.WriteString(strconv.FormatInt(days, 10))
			b.WriteByte('D')
		}
		hasDatePart = true
	}

	if hours > 0 || minutes > 0 || seconds > 0 {
		b.WriteString("T")
		if hours > 0 {
			b.WriteString(strconv.FormatInt(hours, 10))
			b.WriteByte('H')
		}
		if minutes > 0 {
			b.WriteString(strconv.FormatInt(minutes, 10))
			b.WriteByte('M')
		}
		if seconds > 0 {
			b.WriteString(strconv.FormatInt(seconds, 10))
			b.WriteByte('S')
		}
	} else if !hasDatePart {
		b.WriteString("T0S")
	}

	return []byte(`"` + b.String() + `"`), nil
}

// Comment represents a comment on a Yandex Tracker issue.
// Note: Comment.ID is *int (numeric), NOT *string. The LongID field provides
// the string-format ID.
type Comment struct {
	Self        *string       `json:"self,omitempty"`
	ID          *int          `json:"id,omitempty"`
	LongID      *string       `json:"longId,omitempty"`
	Text        *string       `json:"text,omitempty"`
	TextHTML    *string       `json:"textHtml,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	CreatedBy   *User         `json:"createdBy,omitempty"`
	UpdatedBy   *User         `json:"updatedBy,omitempty"`
	CreatedAt   *Timestamp    `json:"createdAt,omitempty"`
	UpdatedAt   *Timestamp    `json:"updatedAt,omitempty"`
	Version     *int          `json:"version,omitempty"`
	Type        *string       `json:"type,omitempty"`
	Transport   *string       `json:"transport,omitempty"`
}

// Attachment represents a file attached to a Yandex Tracker issue or comment.
type Attachment struct {
	Self      *string             `json:"self,omitempty"`
	ID        *string             `json:"id,omitempty"`
	Name      *string             `json:"name,omitempty"`
	Content   *string             `json:"content,omitempty"`
	Thumbnail *string             `json:"thumbnail,omitempty"`
	CreatedBy *User               `json:"createdBy,omitempty"`
	CreatedAt *Timestamp          `json:"createdAt,omitempty"`
	Mimetype  *string             `json:"mimetype,omitempty"`
	Size      *int                `json:"size,omitempty"`
	Metadata  *AttachmentMetadata `json:"metadata,omitempty"`
}

// AttachmentMetadata represents metadata associated with an attachment.
type AttachmentMetadata struct {
	Size *string `json:"size,omitempty"`
}

// ChecklistItem represents an item in a Yandex Tracker issue checklist.
type ChecklistItem struct {
	ID                *string            `json:"id,omitempty"`
	Text              *string            `json:"text,omitempty"`
	TextHTML          *string            `json:"textHtml,omitempty"`
	Checked           *bool              `json:"checked,omitempty"`
	Assignee          *User              `json:"assignee,omitempty"`
	Deadline          *ChecklistDeadline `json:"deadline,omitempty"`
	ChecklistItemType *string            `json:"checklistItemType,omitempty"`
}

// ChecklistDeadline represents the deadline for a checklist item.
type ChecklistDeadline struct {
	Date         *Timestamp `json:"date,omitempty"`
	DeadlineType *string    `json:"deadlineType,omitempty"`
	IsExceeded   *bool      `json:"isExceeded,omitempty"`
}

// Worklog represents a time tracking entry for a Yandex Tracker issue.
type Worklog struct {
	Self      *string    `json:"self,omitempty"`
	ID        *string    `json:"id,omitempty"`
	Version   *string    `json:"version,omitempty"`
	Issue     *Issue     `json:"issue,omitempty"`
	Comment   *string    `json:"comment,omitempty"`
	CreatedBy *User      `json:"createdBy,omitempty"`
	UpdatedBy *User      `json:"updatedBy,omitempty"`
	CreatedAt *Timestamp `json:"createdAt,omitempty"`
	UpdatedAt *Timestamp `json:"updatedAt,omitempty"`
	Start     *Timestamp `json:"start,omitempty"`
	Duration  *Duration  `json:"duration,omitempty"`
}

// QueueVersion represents a version defined in a Yandex Tracker queue.
type QueueVersion struct {
	Self        *string `json:"self,omitempty"`
	ID          *string `json:"id,omitempty"`
	Version     *int    `json:"version,omitempty"`
	Queue       *Queue  `json:"queue,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	StartDate   *string `json:"startDate,omitempty"`
	DueDate     *string `json:"dueDate,omitempty"`
	Released    *bool   `json:"released,omitempty"`
	Archived    *bool   `json:"archived,omitempty"`
}

// QueuePermissions represents the access permissions for a queue.
type QueuePermissions struct {
	Self    *string          `json:"self,omitempty"`
	Version *int             `json:"version,omitempty"`
	Create  *PermissionGroup `json:"create,omitempty"`
	Write   *PermissionGroup `json:"write,omitempty"`
	Read    *PermissionGroup `json:"read,omitempty"`
	Grant   *PermissionGroup `json:"grant,omitempty"`
}

// PermissionGroup represents a group of users, groups, and roles with a permission.
type PermissionGroup struct {
	Self   *string               `json:"self,omitempty"`
	Users  []*PermissionUser     `json:"users,omitempty"`
	Groups []*PermissionGroupRef `json:"groups,omitempty"`
	Roles  []*PermissionRole     `json:"roles,omitempty"`
}

// PermissionUser represents a user in a permission group.
type PermissionUser struct {
	Self        *string `json:"self,omitempty"`
	ID          *string `json:"id,omitempty"`
	Display     *string `json:"display,omitempty"`
	CloudUID    *string `json:"cloudUid,omitempty"`
	PassportUID *string `json:"passportUid,omitempty"`
}

// PermissionGroupRef represents a group reference in a permission group.
type PermissionGroupRef struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// PermissionRole represents a role in a permission group.
type PermissionRole struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// Field represents an issue field in Yandex Tracker (global or local).
type Field struct {
	Self            *string          `json:"self,omitempty"`
	ID              *string          `json:"id,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Key             *string          `json:"key,omitempty"`
	Version         *int             `json:"version,omitempty"`
	Schema          *FieldSchema     `json:"schema,omitempty"`
	Readonly        *bool            `json:"readonly,omitempty"`
	Options         *bool            `json:"options,omitempty"`
	Suggest         *bool            `json:"suggest,omitempty"`
	OptionsProvider *OptionsProvider `json:"optionsProvider,omitempty"`
	QueryProvider   *QueryProvider   `json:"queryProvider,omitempty"`
	Order           *int             `json:"order,omitempty"`
	Category        *FieldCategory   `json:"category,omitempty"`
	Type            *string          `json:"type,omitempty"`
	// Queue is set for local fields (nil for global fields).
	Queue *Queue `json:"queue,omitempty"`
}

// FieldSchema describes the data type of a field.
type FieldSchema struct {
	Type     *string `json:"type,omitempty"`
	Items    *string `json:"items,omitempty"`
	Required *bool   `json:"required,omitempty"`
}

// FieldCategory represents a field category reference.
type FieldCategory struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// OptionsProvider describes the dropdown options configuration for a field.
type OptionsProvider struct {
	Type           *string  `json:"type,omitempty"`
	NeedValidation *bool    `json:"needValidation,omitempty"`
	Values         []string `json:"values,omitempty"`
}

// QueryProvider describes the query provider configuration for a field.
type QueryProvider struct {
	Type *string `json:"type,omitempty"`
}

// FieldName represents a localized field name for create/edit requests.
type FieldName struct {
	EN *string `json:"en,omitempty"`
	RU *string `json:"ru,omitempty"`
}

// Component represents a queue component in Yandex Tracker.
type Component struct {
	Self        *string `json:"self,omitempty"`
	ID          *int    `json:"id,omitempty"`
	Version     *int    `json:"version,omitempty"`
	Name        *string `json:"name,omitempty"`
	Queue       *Queue  `json:"queue,omitempty"`
	Description *string `json:"description,omitempty"`
	Lead        *User   `json:"lead,omitempty"`
	AssignAuto  *bool   `json:"assignAuto,omitempty"`
}

// BulkChange represents the status of an asynchronous bulk change operation.
// Returned by POST /v3/bulkchange/_move, _update, _transition (201 Created)
// and GET /v3/bulkchange/{id} (200 OK).
// Note: BulkChange.ID is *string (hex string like "593cd211ef7e8a33********"), not *int.
type BulkChange struct {
	Self                  *string    `json:"self,omitempty"`
	ID                    *string    `json:"id,omitempty"`
	CreatedBy             *User      `json:"createdBy,omitempty"`
	CreatedAt             *Timestamp `json:"createdAt,omitempty"`
	Status                *string    `json:"status,omitempty"`
	StatusText            *string    `json:"statusText,omitempty"`
	ExecutionChunkPercent *int       `json:"executionChunkPercent,omitempty"`
	ExecutionIssuePercent *int       `json:"executionIssuePercent,omitempty"`
	TotalIssues           *int       `json:"totalIssues,omitempty"`
	TotalCompletedIssues  *int       `json:"totalCompletedIssues,omitempty"`
}

// Board represents an agile board in Yandex Tracker.
type Board struct {
	Self               *string             `json:"self,omitempty"`
	ID                 *int                `json:"id,omitempty"`
	Version            *int                `json:"version,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Columns            []*BoardColumnRef   `json:"columns,omitempty"`
	Filter             map[string]any      `json:"filter,omitempty"`
	OrderBy            *string             `json:"orderBy,omitempty"`
	OrderAsc           *bool               `json:"orderAsc,omitempty"`
	Query              *string             `json:"query,omitempty"`
	UseRanking         *bool               `json:"useRanking,omitempty"`
	Country            *Country            `json:"country,omitempty"`
	DefaultQueue       *Queue              `json:"defaultQueue,omitempty"`
	EstimateBy         *EstimateField      `json:"estimateBy,omitempty"`
	Calendar           *Calendar           `json:"calendar,omitempty"`
	CreatedBy          *User               `json:"createdBy,omitempty"`
	CreatedAt          *Timestamp          `json:"createdAt,omitempty"`
	UpdatedAt          *Timestamp          `json:"updatedAt,omitempty"`
	AutoFilterSettings *AutoFilterSettings `json:"autoFilterSettings,omitempty"`
}

// BoardColumnRef is a lightweight column reference embedded in Board responses.
type BoardColumnRef struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// Country represents a country reference in board responses.
type Country struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// EstimateField represents an estimate field reference in board responses.
type EstimateField struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// Calendar represents a calendar reference in board responses.
type Calendar struct {
	ID *int `json:"id,omitempty"`
}

// AutoFilterSettings represents automatic filter settings on a board.
type AutoFilterSettings struct {
	AddFilterSettings    map[string]any `json:"addFilterSettings,omitempty"`
	RemoveFilterSettings map[string]any `json:"removeFilterSettings,omitempty"`
}

// Column represents a column on an agile board in Yandex Tracker.
// Note: Column.ID is *int (numeric) from the full column response.
// This differs from BoardColumnRef.ID which is *string in embedded board responses.
type Column struct {
	Self     *string   `json:"self,omitempty"`
	ID       *int      `json:"id,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Statuses []*Status `json:"statuses,omitempty"`
}

// Sprint represents an agile sprint in Yandex Tracker.
// Note: Sprint startDate/endDate are plain "YYYY-MM-DD" strings (*string),
// while startDateTime/endDateTime use the full Timestamp format.
type Sprint struct {
	Self          *string    `json:"self,omitempty"`
	ID            *int       `json:"id,omitempty"`
	Version       *int       `json:"version,omitempty"`
	Name          *string    `json:"name,omitempty"`
	Board         *BoardRef  `json:"board,omitempty"`
	Status        *string    `json:"status,omitempty"`
	Archived      *bool      `json:"archived,omitempty"`
	CreatedBy     *User      `json:"createdBy,omitempty"`
	CreatedAt     *Timestamp `json:"createdAt,omitempty"`
	StartDate     *string    `json:"startDate,omitempty"`
	EndDate       *string    `json:"endDate,omitempty"`
	StartDateTime *Timestamp `json:"startDateTime,omitempty"`
	EndDateTime   *Timestamp `json:"endDateTime,omitempty"`
}

// BoardRef is a lightweight board reference embedded in Sprint responses.
type BoardRef struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// Entity represents a Yandex Tracker entity (project, portfolio, or goal).
type Entity struct {
	Self       *string       `json:"self,omitempty"`
	ID         *string       `json:"id,omitempty"`
	Version    *int          `json:"version,omitempty"`
	ShortID    *int          `json:"shortId,omitempty"`
	EntityType *string       `json:"entityType,omitempty"`
	CreatedBy  *User         `json:"createdBy,omitempty"`
	CreatedAt  *Timestamp    `json:"createdAt,omitempty"`
	UpdatedAt  *Timestamp    `json:"updatedAt,omitempty"`
	Fields     *EntityFields `json:"fields,omitempty"`
}

// EntityFields represents the fields of an entity within the "fields" wrapper.
type EntityFields struct {
	Summary              *string          `json:"summary,omitempty"`
	Description          *string          `json:"description,omitempty"`
	Author               *User            `json:"author,omitempty"`
	Lead                 *User            `json:"lead,omitempty"`
	TeamUsers            []*User          `json:"teamUsers,omitempty"`
	Clients              []*User          `json:"clients,omitempty"`
	Followers            []*User          `json:"followers,omitempty"`
	Start                *string          `json:"start,omitempty"`
	End                  *string          `json:"end,omitempty"`
	Tags                 []string         `json:"tags,omitempty"`
	ParentEntity         *ParentEntity    `json:"parentEntity,omitempty"`
	TeamAccess           *bool            `json:"teamAccess,omitempty"`
	EntityStatus         *string          `json:"entityStatus,omitempty"`
	ChecklistItems       []*ChecklistItem `json:"checklistItems,omitempty"`
	IssueQueues          []*Queue         `json:"issueQueues,omitempty"`
	Quarter              []string         `json:"quarter,omitempty"`
	MarkupType           *string          `json:"markupType,omitempty"`
	LastCommentUpdatedAt *string          `json:"lastCommentUpdatedAt,omitempty"`
	LinkedGoalsCount     *int             `json:"linkedGoalsCount,omitempty"`
	LinkedProjectsCount  *int             `json:"linkedProjectsCount,omitempty"`

	// CustomFields holds unknown/type-specific fields not mapped to struct fields.
	CustomFields map[string]any `json:"-"`
}

// entityFieldsKnownKeys lists JSON keys that map to known EntityFields struct fields.
var entityFieldsKnownKeys = []string{
	"summary", "description", "author", "lead", "teamUsers",
	"clients", "followers", "start", "end", "tags",
	"parentEntity", "teamAccess", "entityStatus",
	"checklistItems", "issueQueues", "quarter", "markupType",
	"lastCommentUpdatedAt", "linkedGoalsCount", "linkedProjectsCount",
}

// UnmarshalJSON implements the json.Unmarshaler interface for EntityFields.
// Known fields are decoded into typed struct fields; any unknown keys are
// captured in the CustomFields map.
func (f *EntityFields) UnmarshalJSON(data []byte) error {
	type EntityFieldsAlias EntityFields
	var alias EntityFieldsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*f = EntityFields(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, k := range entityFieldsKnownKeys {
		delete(raw, k)
	}
	if len(raw) > 0 {
		f.CustomFields = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of custom field
			f.CustomFields[k] = val
		}
	}
	return nil
}

// EntitySearchResponse represents the response from entity search.
type EntitySearchResponse struct {
	Hits    *int      `json:"hits,omitempty"`
	Pages   *int      `json:"pages,omitempty"`
	Values  []*Entity `json:"values,omitempty"`
	OrderBy *string   `json:"orderBy,omitempty"`
}

// EntityEventsResponse represents the response from entity event history.
type EntityEventsResponse struct {
	Events  []*EntityEvent `json:"events,omitempty"`
	HasNext *bool          `json:"hasNext,omitempty"`
	HasPrev *bool          `json:"hasPrev,omitempty"`
}

// EntityEvent represents a single event in entity history.
type EntityEvent struct {
	ID        *string              `json:"id,omitempty"`
	Author    *User                `json:"author,omitempty"`
	Date      *Timestamp           `json:"date,omitempty"`
	Transport *string              `json:"transport,omitempty"`
	Display   *string              `json:"display,omitempty"`
	Changes   []*EntityEventChange `json:"changes,omitempty"`
}

// EntityEventChange represents a single field change within an event.
type EntityEventChange struct {
	Diff  *string           `json:"diff,omitempty"`
	Field *EntityEventField `json:"field,omitempty"`
}

// EntityEventField identifies the field that changed.
type EntityEventField struct {
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// EntityAccess represents the access settings for an entity.
type EntityAccess struct {
	ACL               *EntityACL            `json:"acl,omitempty"`
	PermissionSources []*EntityRef          `json:"permissionSources,omitempty"`
	ParentEntities    *EntityParentEntities `json:"parentEntities,omitempty"`
}

// EntityACL represents the access control list for an entity.
type EntityACL struct {
	Read  *EntityACLEntry `json:"READ,omitempty"`
	Write *EntityACLEntry `json:"WRITE,omitempty"`
	Grant *EntityACLEntry `json:"GRANT,omitempty"`
}

// EntityACLEntry represents a permission entry with users, groups, and roles.
type EntityACLEntry struct {
	Users  []*User               `json:"users,omitempty"`
	Groups []*PermissionGroupRef `json:"groups,omitempty"`
	Roles  []string              `json:"roles,omitempty"`
}

// EntityRef is a lightweight entity reference used in access settings.
type EntityRef struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// EntityParentEntities represents the parent entity hierarchy in access settings.
type EntityParentEntities struct {
	Primary   *EntityRef   `json:"primary,omitempty"`
	Secondary []*EntityRef `json:"secondary,omitempty"`
}

// ParentEntity represents the parent entity reference in entity fields.
type ParentEntity struct {
	Primary   *string  `json:"primary,omitempty"`
	Secondary []string `json:"secondary,omitempty"`
}

// EntityLinkResponse represents a link between entities as returned by the API.
// Unlike IssueLink which contains an embedded Issue object, entity links
// return the link type and requested field values of the linked entity.
type EntityLinkResponse struct {
	Type            *string        `json:"type,omitempty"`
	LinkFieldValues map[string]any `json:"linkFieldValues,omitempty"`
}

// Trigger represents a queue trigger in Yandex Tracker.
// Triggers automate actions when specific conditions are met on issues.
type Trigger struct {
	ID         *string                `json:"id,omitempty"`
	Self       *string                `json:"self,omitempty"`
	Queue      *Queue                 `json:"queue,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Order      *string                `json:"order,omitempty"`
	Actions    []*AutomationAction    `json:"actions,omitempty"`
	Conditions []*AutomationCondition `json:"conditions,omitempty"`
	Version    *int                   `json:"version,omitempty"`
	Active     *bool                  `json:"active,omitempty"`
}

// AutoAction represents a queue auto-action in Yandex Tracker.
// Auto-actions are periodic automation rules that run on a schedule,
// filtering issues by filter/query and applying actions.
type AutoAction struct {
	ID                   *string             `json:"id,omitempty"`
	Self                 *string             `json:"self,omitempty"`
	Queue                *Queue              `json:"queue,omitempty"`
	Name                 *string             `json:"name,omitempty"`
	Version              *int                `json:"version,omitempty"`
	Active               *bool               `json:"active,omitempty"`
	Created              *Timestamp          `json:"created,omitempty"`
	Updated              *Timestamp          `json:"updated,omitempty"`
	Filter               any                 `json:"filter,omitempty"`
	Query                *string             `json:"query,omitempty"`
	Actions              []*AutomationAction `json:"actions,omitempty"`
	EnableNotifications  *bool               `json:"enableNotifications,omitempty"`
	LastLaunch           *Timestamp          `json:"lastLaunch,omitempty"`
	TotalIssuesProcessed *int                `json:"totalIssuesProcessed,omitempty"`
	IntervalMillis       *int64              `json:"intervalMillis,omitempty"`
	Calendar             *AutoActionCalendar `json:"calendar,omitempty"`
}

// AutoActionCalendar represents a calendar schedule reference for an auto-action.
type AutoActionCalendar struct {
	ID *int `json:"id,omitempty"`
}

// AutomationAction represents an action in a trigger or auto-action.
// Known fields (Type, ID, Status) are typed struct fields. All other
// action-type-specific fields are captured in the Parameters map.
// The custom MarshalJSON/UnmarshalJSON methods flatten Parameters into
// the top-level JSON object alongside the known fields.
type AutomationAction struct {
	Type       *string        `json:"-"`
	ID         *string        `json:"-"`
	Status     *Status        `json:"-"`
	Parameters map[string]any `json:"-"`
}

// automationActionKnownKeys lists the JSON keys that map to known
// AutomationAction struct fields.
var automationActionKnownKeys = []string{"type", "id", "status"}

// UnmarshalJSON implements the json.Unmarshaler interface for AutomationAction.
// It decodes the known fields (type, id, status) into struct fields and stores
// any remaining keys in the Parameters map.
func (a *AutomationAction) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["type"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			a.Type = &s
		}
	}
	if v, ok := raw["id"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			a.ID = &s
		}
	}
	if v, ok := raw["status"]; ok {
		st := new(Status)
		if err := json.Unmarshal(v, st); err == nil {
			a.Status = st
		}
	}

	// Remaining keys go into Parameters.
	for _, k := range automationActionKnownKeys {
		delete(raw, k)
	}
	if len(raw) > 0 {
		a.Parameters = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of parameter
			a.Parameters[k] = val
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for AutomationAction.
// It builds a flat JSON object containing both the known fields and
// all entries from the Parameters map.
func (a AutomationAction) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	// Copy overflow parameters first.
	for k, v := range a.Parameters {
		m[k] = v
	}

	// Overlay known fields (they take precedence).
	if a.Type != nil {
		m["type"] = *a.Type
	}
	if a.ID != nil {
		m["id"] = *a.ID
	}
	if a.Status != nil {
		m["status"] = a.Status
	}

	return json.Marshal(m)
}

// AutomationCondition represents a condition in a trigger or auto-action.
// Known fields (Type, Conditions) are typed struct fields. All other
// condition-type-specific fields are captured in the Parameters map.
// The custom MarshalJSON/UnmarshalJSON methods flatten Parameters into
// the top-level JSON object alongside the known fields.
type AutomationCondition struct {
	Type       *string                `json:"-"`
	Conditions []*AutomationCondition `json:"-"`
	Parameters map[string]any         `json:"-"`
}

// automationConditionKnownKeys lists the JSON keys that map to known
// AutomationCondition struct fields.
var automationConditionKnownKeys = []string{"type", "conditions"}

// UnmarshalJSON implements the json.Unmarshaler interface for AutomationCondition.
// It decodes the known fields (type, conditions) into struct fields and stores
// any remaining keys in the Parameters map.
func (c *AutomationCondition) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["type"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			c.Type = &s
		}
	}
	if v, ok := raw["conditions"]; ok {
		var conds []*AutomationCondition
		if err := json.Unmarshal(v, &conds); err == nil {
			c.Conditions = conds
		}
	}

	// Remaining keys go into Parameters.
	for _, k := range automationConditionKnownKeys {
		delete(raw, k)
	}
	if len(raw) > 0 {
		c.Parameters = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of parameter
			c.Parameters[k] = val
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for AutomationCondition.
// It builds a flat JSON object containing both the known fields and
// all entries from the Parameters map.
func (c AutomationCondition) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	// Copy overflow parameters first.
	for k, v := range c.Parameters {
		m[k] = v
	}

	// Overlay known fields (they take precedence).
	if c.Type != nil {
		m["type"] = *c.Type
	}
	if c.Conditions != nil {
		m["conditions"] = c.Conditions
	}

	return json.Marshal(m)
}

// Macro represents a macro in Yandex Tracker.
// Macros allow batch operations on issues triggered manually by users.
type Macro struct {
	ID          *int                `json:"id,omitempty"`
	Self        *string             `json:"self,omitempty"`
	Queue       *Queue              `json:"queue,omitempty"`
	Name        *string             `json:"name,omitempty"`
	Body        *string             `json:"body,omitempty"`
	IssueUpdate []*MacroIssueUpdate `json:"issueUpdate,omitempty"`
}

// MacroIssueUpdate represents a field update entry in a macro's issueUpdate array.
type MacroIssueUpdate struct {
	Field  *MacroIssueUpdateField `json:"field,omitempty"`
	Update any                    `json:"update,omitempty"`
}

// MacroIssueUpdateField represents the field metadata within a macro issue update.
type MacroIssueUpdateField struct {
	Self    *string `json:"self,omitempty"`
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// Filter represents a saved filter in Yandex Tracker.
type Filter struct {
	Self    *string        `json:"self,omitempty"`
	ID      *int           `json:"id,omitempty"`
	Name    *string        `json:"name,omitempty"`
	Filter  map[string]any `json:"filter,omitempty"`
	Query   *string        `json:"query,omitempty"`
	Sorts   any            `json:"sorts,omitempty"`
	GroupBy *string        `json:"groupBy,omitempty"`
}

// ExternalApplication represents an external application registered in Yandex Tracker.
type ExternalApplication struct {
	Self *string `json:"self,omitempty"`
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ExternalLink represents an external link on an issue in Yandex Tracker.
type ExternalLink struct {
	Self      *string             `json:"self,omitempty"`
	ID        *int                `json:"id,omitempty"`
	Type      *ExternalLinkType   `json:"type,omitempty"`
	Direction *string             `json:"direction,omitempty"`
	Object    *ExternalLinkObject `json:"object,omitempty"`
	CreatedBy *User               `json:"createdBy,omitempty"`
	UpdatedBy *User               `json:"updatedBy,omitempty"`
	CreatedAt *Timestamp          `json:"createdAt,omitempty"`
	UpdatedAt *Timestamp          `json:"updatedAt,omitempty"`
}

// ExternalLinkType represents the type information of an external link.
type ExternalLinkType struct {
	Self *string `json:"self,omitempty"`
	ID   *string `json:"id,omitempty"`
	Key  *string `json:"key,omitempty"`
}

// ExternalLinkObject represents the external object referenced by an external link.
type ExternalLinkObject struct {
	Self        *string              `json:"self,omitempty"`
	ID          *string              `json:"id,omitempty"`
	Key         *string              `json:"key,omitempty"`
	Application *ExternalApplication `json:"application,omitempty"`
}

// Dashboard represents a Yandex Tracker dashboard.
type Dashboard struct {
	Self      *string    `json:"self,omitempty"`
	ID        *int       `json:"id,omitempty"`
	Version   *int       `json:"version,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Layout    *string    `json:"layout,omitempty"`
	Owner     *User      `json:"owner,omitempty"`
	CreatedBy *User      `json:"createdBy,omitempty"`
	CreatedAt *Timestamp `json:"createdAt,omitempty"`
}

// Widget represents a widget on a Yandex Tracker dashboard.
// Known fields are extracted into struct fields; any additional
// type-specific fields are stored in the Parameters map.
type Widget struct {
	Self        *string        `json:"-"`
	ID          *int           `json:"-"`
	Version     *int           `json:"-"`
	Description *string        `json:"-"`
	CreatedBy   *User          `json:"-"`
	Dashboard   *DashboardRef  `json:"-"`
	DatasetInfo *DatasetInfo   `json:"-"`
	Parameters  map[string]any `json:"-"`
}

// widgetKnownKeys lists the JSON keys that map to known
// Widget struct fields.
var widgetKnownKeys = []string{"self", "id", "version", "description", "createdBy", "dashboard", "datasetInfo"}

// UnmarshalJSON implements the json.Unmarshaler interface for Widget.
// It decodes the known fields into struct fields and stores
// any remaining keys in the Parameters map.
func (w *Widget) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["self"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			w.Self = &s
		}
	}
	if v, ok := raw["id"]; ok {
		var n int
		if err := json.Unmarshal(v, &n); err == nil {
			w.ID = &n
		}
	}
	if v, ok := raw["version"]; ok {
		var n int
		if err := json.Unmarshal(v, &n); err == nil {
			w.Version = &n
		}
	}
	if v, ok := raw["description"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			w.Description = &s
		}
	}
	if v, ok := raw["createdBy"]; ok {
		var u User
		if err := json.Unmarshal(v, &u); err == nil {
			w.CreatedBy = &u
		}
	}
	if v, ok := raw["dashboard"]; ok {
		var d DashboardRef
		if err := json.Unmarshal(v, &d); err == nil {
			w.Dashboard = &d
		}
	}
	if v, ok := raw["datasetInfo"]; ok {
		var d DatasetInfo
		if err := json.Unmarshal(v, &d); err == nil {
			w.DatasetInfo = &d
		}
	}

	// Remaining keys go into Parameters.
	for _, k := range widgetKnownKeys {
		delete(raw, k)
	}
	if len(raw) > 0 {
		w.Parameters = make(map[string]any, len(raw))
		for k, v := range raw {
			var val any
			_ = json.Unmarshal(v, &val) // best-effort decode of parameter
			w.Parameters[k] = val
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for Widget.
// It builds a flat JSON object containing both the known fields and
// all entries from the Parameters map.
func (w Widget) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	// Copy overflow parameters first.
	for k, v := range w.Parameters {
		m[k] = v
	}

	// Overlay known fields (they take precedence).
	if w.Self != nil {
		m["self"] = *w.Self
	}
	if w.ID != nil {
		m["id"] = *w.ID
	}
	if w.Version != nil {
		m["version"] = *w.Version
	}
	if w.Description != nil {
		m["description"] = *w.Description
	}
	if w.CreatedBy != nil {
		m["createdBy"] = w.CreatedBy
	}
	if w.Dashboard != nil {
		m["dashboard"] = w.Dashboard
	}
	if w.DatasetInfo != nil {
		m["datasetInfo"] = w.DatasetInfo
	}

	return json.Marshal(m)
}

// DashboardRef is a lightweight dashboard reference in widget responses.
type DashboardRef struct {
	ID      *string `json:"id,omitempty"`
	Display *string `json:"display,omitempty"`
}

// DatasetInfo contains widget calculation metadata.
type DatasetInfo struct {
	Status         *string `json:"status,omitempty"`
	BuildStartedAt *string `json:"buildStartedAt,omitempty"`
}
