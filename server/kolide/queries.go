package kolide

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
)

type QueryStore interface {
	// Query methods
	NewQuery(query *Query) (*Query, error)
	SaveQuery(query *Query) error
	DeleteQuery(query *Query) error
	Query(id uint) (*Query, error)
	ListQueries(opt ListOptions) ([]*Query, error)

	// NewDistributedQueryCampaign creates a new distributed query campaign
	NewDistributedQueryCampaign(camp DistributedQueryCampaign) (DistributedQueryCampaign, error)
	// SaveDistributedQueryCampaign updates an existing distributed query
	// campaign
	SaveDistributedQueryCampaign(camp DistributedQueryCampaign) error
	// NewDistributedQueryCampaignTarget adds a new target to an existing
	// distributed query campaign
	NewDistributedQueryCampaignTarget(target DistributedQueryCampaignTarget) (DistributedQueryCampaignTarget, error)
	// NewDistributedQueryCampaignExecution records a new execution for a
	// distributed query campaign
	NewDistributedQueryExecution(exec DistributedQueryExecution) (DistributedQueryExecution, error)
}

type QueryService interface {
	ListQueries(ctx context.Context, opt ListOptions) ([]*Query, error)
	GetQuery(ctx context.Context, id uint) (*Query, error)
	NewQuery(ctx context.Context, p QueryPayload) (*Query, error)
	ModifyQuery(ctx context.Context, id uint, p QueryPayload) (*Query, error)
	DeleteQuery(ctx context.Context, id uint) error
}

type QueryPayload struct {
	Name         *string
	Description  *string
	Query        *string
	Interval     *uint
	Snapshot     *bool
	Differential *bool
	Platform     *string
	Version      *string
}

type PackPayload struct {
	Name     *string
	Platform *string
}

type Query struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	Name         string    `json:"name" gorm:"not null;unique_index:idx_query_unique_name"`
	Description  string    `json:"description"`
	Query        string    `json:"query" gorm:"not null"`
	Interval     uint      `json:"interval"`
	Snapshot     bool      `json:"snapshot"`
	Differential bool      `json:"differential"`
	Platform     string    `json:"platform"`
	Version      string    `json:"version"`
}

type DistributedQueryStatus int

const (
	QueryRunning  DistributedQueryStatus = iota
	QueryComplete DistributedQueryStatus = iota
	QueryError    DistributedQueryStatus = iota
)

type DistributedQueryCampaign struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	QueryID     uint
	MaxDuration time.Duration
	Status      DistributedQueryStatus
	UserID      uint
}

type DistributedQueryCampaignTarget struct {
	ID                         uint `gorm:"primary_key"`
	Type                       TargetType
	DistributedQueryCampaignID uint `gorm:"index:idx_dqct_dqc_id"`
	TargetID                   uint
}

type DistributedQueryExecutionStatus int

const (
	ExecutionWaiting DistributedQueryExecutionStatus = iota
	ExecutionRequested
	ExecutionSucceeded
	ExecutionFailed
)

type DistributedQueryResult struct {
	DistributedQueryCampaignID uint            `json:"distributed_query_execution_id"`
	Host                       Host            `json:"host"`
	ResultJSON                 json.RawMessage `json:"result_json"`
}

type DistributedQueryExecution struct {
	ID                         uint `gorm:"primary_key"`
	HostID                     uint // unique index added in migrate
	DistributedQueryCampaignID uint // unique index added in migrate
	Status                     DistributedQueryExecutionStatus
	Error                      string `gorm:"size:1024"`
	ExecutionDuration          time.Duration
}

type Option struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"not null;unique_index:idx_option_unique_key"`
	Value     string `gorm:"not null"`
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad DecoratorType = iota
	DecoratorAlways
	DecoratorInterval
)

type Decorator struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType `gorm:"not null"`
	Interval  int
	Query     string
}
