package service

import (
	"context"
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Enroll Agent
////////////////////////////////////////////////////////////////////////////////

type enrollAgentRequest struct {
	EnrollSecret   string                         `json:"enroll_secret"`
	HostIdentifier string                         `json:"host_identifier"`
	HostDetails    map[string](map[string]string) `json:"host_details"`
}

type enrollAgentResponse struct {
	NodeKey string `json:"node_key,omitempty"`
	Err     error  `json:"error,omitempty"`
}

func (r enrollAgentResponse) error() error { return r.Err }

func makeEnrollAgentEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(enrollAgentRequest)
		nodeKey, err := svc.EnrollAgent(ctx, req.EnrollSecret, req.HostIdentifier, req.HostDetails)
		if err != nil {
			return enrollAgentResponse{Err: err}, nil
		}
		return enrollAgentResponse{NodeKey: nodeKey}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Client Config
////////////////////////////////////////////////////////////////////////////////

type getClientConfigRequest struct {
	NodeKey string `json:"node_key"`
}

type getClientConfigResponse struct {
	Config map[string]interface{}
	Err    error `json:"error,omitempty"`
}

func (r getClientConfigResponse) error() error { return r.Err }

func makeGetClientConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.GetClientConfig(ctx)
		if err != nil {
			return getClientConfigResponse{Err: err}, nil
		}

		// We return the config here explicitly because osquery exepects the
		// response for configs to be at the top-level of the JSON response
		return config, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Distributed Queries
////////////////////////////////////////////////////////////////////////////////

type getDistributedQueriesRequest struct {
	NodeKey string `json:"node_key"`
}

type getDistributedQueriesResponse struct {
	Queries    map[string]string `json:"queries"`
	Accelerate uint              `json:"accelerate,omitempty"`
	Err        error             `json:"error,omitempty"`
}

func (r getDistributedQueriesResponse) error() error { return r.Err }

func makeGetDistributedQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		queries, accelerate, err := svc.GetDistributedQueries(ctx)
		if err != nil {
			return getDistributedQueriesResponse{Err: err}, nil
		}
		return getDistributedQueriesResponse{Queries: queries, Accelerate: accelerate}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Write Distributed Query Results
////////////////////////////////////////////////////////////////////////////////

type submitDistributedQueryResultsRequest struct {
	NodeKey  string                                `json:"node_key"`
	Results  kolide.OsqueryDistributedQueryResults `json:"queries"`
	Statuses map[string]kolide.OsqueryStatus       `json:"statuses"`
}

type submitDistributedQueryResultsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitDistributedQueryResultsResponse) error() error { return r.Err }

func makeSubmitDistributedQueryResultsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitDistributedQueryResultsRequest)
		err := svc.SubmitDistributedQueryResults(ctx, req.Results, req.Statuses)
		if err != nil {
			return submitDistributedQueryResultsResponse{Err: err}, nil
		}
		return submitDistributedQueryResultsResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Submit Logs
////////////////////////////////////////////////////////////////////////////////

type submitLogsRequest struct {
	NodeKey string          `json:"node_key"`
	LogType string          `json:"log_type"`
	Data    json.RawMessage `json:"data"`
}

type submitLogsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitLogsResponse) error() error { return r.Err }

func makeSubmitLogsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitLogsRequest)

		var err error
		switch req.LogType {
		case "status":
			var statuses []json.RawMessage
			if err := json.Unmarshal(req.Data, &statuses); err != nil {
				err = osqueryError{message: "unmarshalling status logs: " + err.Error()}
				break
			}

			err = svc.SubmitStatusLogs(ctx, statuses)
			if err != nil {
				break
			}

		case "result":
			var results []json.RawMessage
			if err := json.Unmarshal(req.Data, &results); err != nil {
				err = osqueryError{message: "unmarshalling result logs: " + err.Error()}
				break
			}
			err = svc.SubmitResultLogs(ctx, results)
			if err != nil {
				break
			}

		default:
			err = osqueryError{message: "unknown log type: " + req.LogType}
		}

		return submitLogsResponse{Err: err}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Begin File Carve
////////////////////////////////////////////////////////////////////////////////

type carveBeginRequest struct {
	NodeKey    string `json:"node_key"`
	BlockCount int    `json:"block_count"`
	BlockSize  int    `json:"block_size"`
	CarveSize  int    `json:"carve_size"`
	CarveId    string `json:"carve_id"`
	RequestId  string `json:"request_id"`
}

type carveBeginResponse struct {
	SessionId string `json:"session_id"`
	Success   bool   `json:"success,omitempty"`
	Err       error  `json:"error,omitempty"`
}

func (r carveBeginResponse) error() error { return r.Err }

func makeCarveBeginEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(carveBeginRequest)
		_ = req

		spew.Dump(req)
		// TODO call some method
		id := "foobar"
		var err error

		if err != nil {
			return carveBeginResponse{Err: err}, nil
		}

		return carveBeginResponse{SessionId: id, Success: true}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Receive Block for File Carve
////////////////////////////////////////////////////////////////////////////////

type carveBlockRequest struct {
	NodeKey   string `json:"node_key"`
	BlockId   int    `json:"block_id"`
	SessionId string `json:"session_id"`
	RequestId string `json:"request_id"`
	Data      string `json:"data"`
}

type carveBlockResponse struct {
	Success bool  `json:"success,omitempty"`
	Err     error `json:"error,omitempty"`
}

func (r carveBlockResponse) error() error { return r.Err }

func makeCarveBlockEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(carveBlockRequest)
		_ = req
		spew.Dump(req)

		// TODO call method
		var err error

		if err != nil {
			return carveBlockResponse{Err: err}, nil
		}

		return carveBlockResponse{Success: true}, nil
	}
}
