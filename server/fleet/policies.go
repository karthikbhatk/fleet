package fleet

// PolicyPayload holds data for policy creation.
type PolicyPayload struct {
	// QueryID allows creating a policy from an existing query.
	//
	// Using QueryID is the old way of creating policies.
	// Use Query, Name and Description instead.
	QueryID uint
	// Name is the name of the policy (ignored if QueryID != 0).
	Name string
	// Query is the policy query (ignored if QueryID != 0).
	Query string
	// Description is the policy description text (ignored if QueryID != 0).
	Description string
	// Resolution indicate the steps needed to solve a failing policy.
	Resolution string
}

// ModifyPolicyPayload holds data for policy modification.
type ModifyPolicyPayload struct {
	// Name is the name of the policy.
	Name *string `json:"name"`
	// Query is the policy query.
	Query *string `json:"query"`
	// Description is the policy description text.
	Description *string `json:"description"`
	// Resolution indicate the steps needed to solve a failing policy.
	Resolution *string `json:"resolution"`
}

// Policy is a fleet's policy query.
//
// TODO(lucas): Add AuthorEmail mimicking Query.
type Policy struct {
	ID uint `json:"id"`
	// Name is the name of the policy query.
	// NOTE(lucas): To not break clients (UI), I'm not changing this to json:"name" (#2595).
	Name        string `json:"query_name" db:"name"`
	Query       string `json:"query" db:"query"`
	Description string `json:"description" db:"description"`
	AuthorID    *uint  `json:"author_id" db:"author_id"`
	// AuthorName is retrieved with a join to the users table in the MySQL backend (using AuthorID).
	AuthorName       string  `json:"author_name" db:"author_name"`
	PassingHostCount uint    `json:"passing_host_count" db:"passing_host_count"`
	FailingHostCount uint    `json:"failing_host_count" db:"failing_host_count"`
	TeamID           *uint   `json:"team_id" db:"team_id"`
	Resolution       *string `json:"resolution,omitempty" db:"resolution"`

	TeamIDX uint `json:"-" db:"team_id_x"`

	UpdateCreateTimestamps
}

func (p Policy) AuthzType() string {
	return "policy"
}

const (
	PolicyKind = "policy"
)

type HostPolicy struct {
	ID uint `json:"id" db:"id"`
	// Name is the name of the policy query.
	// NOTE(lucas): To not break clients (UI), I'm not changing this to json:"name" (#2595).
	Name  string `json:"query_name" db:"name"`
	Query string `json:"query" db:"query"`
	// Description is the policy description.
	// NOTE(lucas): To not break clients (UI), I'm not changing this to json:"description" (#2595).
	Description string `json:"query_description" db:"description"`
	AuthorID    *uint  `json:"author_id" db:"author_id"`
	// AuthorName is retrieved with a join to the users table in the MySQL backend (using AuthorID).
	AuthorName string `json:"author_name" db:"author_name"`
	Response   string `json:"response" db:"response"`
	Resolution string `json:"resolution" db:"resolution"`

	TeamIDX uint `json:"-" db:"team_id_x"`
}

type PolicySpec struct {
	Name        string `json:"name"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Resolution  string `json:"resolution,omitempty"`
	Team        string `json:"team,omitempty"`
}
