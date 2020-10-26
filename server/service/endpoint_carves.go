package service

import (
	"context"
	"io"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Begin File Carve
////////////////////////////////////////////////////////////////////////////////

type carveBeginRequest struct {
	NodeKey    string `json:"node_key"`
	BlockCount int64  `json:"block_count"`
	BlockSize  int64  `json:"block_size"`
	CarveSize  int64  `json:"carve_size"`
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

		payload := kolide.CarveBeginPayload{
			BlockCount: req.BlockCount,
			BlockSize:  req.BlockSize,
			CarveSize:  req.CarveSize,
			CarveId:    req.CarveId,
			RequestId:  req.RequestId,
		}

		carve, err := svc.CarveBegin(ctx, payload)
		if err != nil {
			return carveBeginResponse{Err: err}, nil
		}

		return carveBeginResponse{SessionId: carve.SessionId, Success: true}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Receive Block for File Carve
////////////////////////////////////////////////////////////////////////////////

type carveBlockRequest struct {
	NodeKey   string `json:"node_key"`
	BlockId   int64  `json:"block_id"`
	SessionId string `json:"session_id"`
	RequestId string `json:"request_id"`
	Data      []byte `json:"data"`
}

type carveBlockResponse struct {
	Success bool  `json:"success,omitempty"`
	Err     error `json:"error,omitempty"`
}

func (r carveBlockResponse) error() error { return r.Err }

func makeCarveBlockEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(carveBlockRequest)

		payload := kolide.CarveBlockPayload{
			SessionId: req.SessionId,
			RequestId: req.RequestId,
			BlockId:   req.BlockId,
			Data:      req.Data,
		}

		err := svc.CarveBlock(ctx, payload)
		if err != nil {
			return carveBlockResponse{Err: err}, nil
		}

		return carveBlockResponse{Success: true}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// List Carves
////////////////////////////////////////////////////////////////////////////////

type carveResponse struct {
	kolide.CarveMetadata
}

type listCarvesRequest struct {
	ListOptions kolide.ListOptions
}

type listCarvesResponse struct {
	Carves []carveResponse `json:"carves"`
	Err    error           `json:"error,omitempty"`
}

func (r listCarvesResponse) error() error { return r.Err }

func makeListCarvesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listCarvesRequest)
		carves, err := svc.ListCarves(ctx, req.ListOptions)
		if err != nil {
			return listCarvesResponse{Err: err}, nil
		}

		resp := listCarvesResponse{}
		for _, carve := range carves {
			resp.Carves = append(resp.Carves, carveResponse{*carve})
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Download Carve
////////////////////////////////////////////////////////////////////////////////

type downloadCarveRequest struct {
	Name string
}

type downloadCarveResponse struct {
	Reader io.Reader `json:"-"`
	Err    error     `json:"error,omitempty"`
}

func (r downloadCarveResponse) error() error { return r.Err }

func makeDownloadCarveEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(downloadCarveRequest)
		reader, err := svc.GetCarveReader(ctx, req.Name)
		if err != nil {
			return downloadCarveResponse{Err: err}, nil
		}

		return downloadCarveResponse{Reader: reader}, nil
	}
}
